// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package web

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/chorse-dev/cdash-proxy/ctestxml"
	"github.com/chorse-dev/cdash-proxy/gcovtar"
)

func Put(w http.ResponseWriter, r *http.Request, hf HandlerFunc) {
	if r.FormValue("type") == "GcovTar" {
		buildID := r.FormValue("buildid")
		job, err := gcovtar.Parse(r.Body, buildID)
		if err == nil {
			err = hf(r.Context(), job)
		}
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			return
		}
		fmt.Fprintf(w, `{"status":0}`)
		return
	}

	if filepath.Ext(r.FormValue("FileName")) == ".xml" {
		buildID, err := handleCTestXML(r, hf)

		sendResponseXML(w, buildID, err)
		return
	}

	http.NotFound(w, r)
}

func handleCTestXML(r *http.Request, hf HandlerFunc) (string, error) {
	job, err := ctestxml.Parse(r.Body, r.FormValue("project"))
	if err != nil {
		return "", err
	}

	if err := hf(r.Context(), job); err != nil {
		return "", err
	}

	return job.JobID, nil
}

func sendResponseXML(w http.ResponseWriter, buildID string, err error) {
	var resp *ctestxml.Response
	if err == nil {
		resp = &ctestxml.Response{
			Status:  "OK",
			BuildID: buildID,
		}
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		resp = &ctestxml.Response{
			Status:  "ERROR",
			Message: err.Error(),
		}
	}
	if err = xml.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
