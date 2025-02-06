// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/purpleKarrot/cdash-proxy/model"
	"github.com/purpleKarrot/cdash-proxy/web"
)

func print(_ context.Context, job *model.Job) error {
	jobJSON, err := json.Marshal(job)
	if err != nil {
		return err
	}

	fmt.Println(string(jobJSON))
	return nil
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(web.Serve(print))))
}
