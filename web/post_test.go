// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPost(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	w := httptest.NewRecorder()
	Post(w, r)
	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}
