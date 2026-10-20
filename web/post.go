// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package web

import (
	"encoding/json"
	"net/http"

	"github.com/chorse-dev/cdash-proxy/ctestxml"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid upload form", http.StatusBadRequest)
		return
	}
	jobID := ctestxml.GenerateJobID(
		r.Form.Get("project"),
		r.Form.Get("site"),
		r.Form.Get("stamp"),
		r.Form.Get("build"),
	)

	type uploadRequest struct {
		Status       int    `json:"status"`
		DataFilesMD5 []int  `json:"datafilesmd5"`
		JobID        string `json:"buildid"`
	}

	json.NewEncoder(w).Encode(uploadRequest{
		Status:       0,
		DataFilesMD5: []int{0},
		JobID:        jobID,
	})
}
