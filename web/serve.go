// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package web

import (
	"context"
	"net/http"

	"github.com/chorse-dev/cdash-proxy/model"
)

type HandlerFunc func(ctx context.Context, job *model.Job) error

func Serve(hf HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			Put(w, r, hf)
		case "POST":
			Post(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
		}
	}
}
