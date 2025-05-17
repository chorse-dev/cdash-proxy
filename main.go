// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chorse-dev/cdash-proxy/model"
	"github.com/chorse-dev/cdash-proxy/web"
)

func print(_ context.Context, job *model.Job) error {
	jobJSON, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jobJSON))
	return nil
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(web.Serve(print))))
}
