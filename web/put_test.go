// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package web

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/purpleKarrot/cdash-proxy/model"
)

var serveOK = Serve(func(ctx context.Context, job *model.Job) error {
	return nil
})

var serveError = Serve(func(ctx context.Context, job *model.Job) error {
	return errors.New("test error")
})

func TestPutGcovTarOK(t *testing.T) {
	file, _ := os.Open("../gcovtar/testdata/gcov.tbz2")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/submit?type=GcovTar", file)
	serveOK(w, r)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(body) != `{"status":0}` {
		t.Errorf("unexpected body: %s", string(body))
	}
}

func TestPutGcovTarError(t *testing.T) {
	file, _ := os.Open("../gcovtar/testdata/gcov.tbz2")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/submit?type=GcovTar", file)
	serveError(w, r)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(body) != `{"error":"test error"}` {
		t.Errorf("unexpected body: %s", string(body))
	}
}

func TestPutXmlOK(t *testing.T) {
	file, _ := os.Open("../ctestxml/testdata/Configure.xml")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/submit?project=Example&FileName=Configure.xml", file)
	serveOK(w, r)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(body) != `<cdash><status>OK</status><buildId>4e5a4b59fc4badd8ec47227aa4514ba1</buildId></cdash>` {
		t.Errorf("unexpected body: %s", string(body))
	}
}

func TestPutXmlError(t *testing.T) {
	file, _ := os.Open("../ctestxml/testdata/Configure.xml")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/submit?project=Example&FileName=Configure.xml", file)
	serveError(w, r)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(body) != `<cdash><status>ERROR</status><message>test error</message></cdash>` {
		t.Errorf("unexpected body: %s", string(body))
	}
}
