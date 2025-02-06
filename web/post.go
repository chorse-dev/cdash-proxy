// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package web

import (
	"encoding/json"
	"net/http"

	"github.com/purpleKarrot/cdash-proxy/ctestxml"
)

func Post(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	jobID := ctestxml.GenerateJobID(
		query.Get("project"),
		query.Get("site"),
		query.Get("stamp"),
		query.Get("build"),
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
