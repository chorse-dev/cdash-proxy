// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package util

import (
	"testing"
)

func TestDecodeOK(t *testing.T) {
	input := `eJxTUFBQKEktLlHILy0pKC3hAnIVcvOLUvHxuQC1kg/C====`
	expected := "    test output\n    more output\n    more output\n    \n"
	buf, err := Decode(input, "base64", "gzip")
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(buf) != expected {
		t.Errorf("decoded: '%v' , '%v'", string(buf), buf)
	}
}

func TestDecodeFail(t *testing.T) {
	input := `exTUFBQKEktLlHILy0pKC3hAnIVcvOLUvHxuQC1kg/C====`
	_, err := Decode(input, "base64", "gzip")
	if err == nil || err.Error() != "zlib: invalid header" {
		t.Errorf("expected error, got '%v'", err)
	}
}

func TestDecodeTar(t *testing.T) {
	input := `H4sIADBZ42YAA+1V226bQBDNM18xIn1IpQSwS+yqSq1aQBKrJrFskiiKLLSBIax` +
		`iWLQsda3KH9Tf6Jd1wcS5WE2sRulFyryM5sI5c3ZYsFxyjX2ai1wTX8XGi5hhGC3TBOk` +
		`b7V3jrq+s2X7XgIbZahutXdNsGWA0dptmewOMlxnnvhW5IFyOEmYR0gj5r/qeqi+0wNL` +
		`/J7YJFstmnF7FAn58h6bRNMEmKcUJDBaCYS+s4k/1AewkhE60EDvKprIJXkxzmDJ+DdJ` +
		`HHFEDOGcFBCQFjqF8sTi9LAQCFUDSUGccEhbSaFYmijSU+CJGCSSQJzmwqAzBZnAWEyH` +
		`REfaL4LqCPCOpzDAYFJcTGkCfBpjmuA2nyHPKUmhuSxiSQ1bW8xhDuJzBiCRwyAL8Qri` +
		`cbIQIsRDZB12fTqfaVETZREtR6BBVg3GEEIXUl2uKkqPY6rndA8e3e0NQ33yz3O5nx7d` +
		`OhkPnyPNHxydDq6rN1beKQsLQF5iLraOu60BaJMhp4CdI8oJjgqlQAKxj1+0e2bCEWsR` +
		`z2HEAg5jBxYXsAihxgBUiK0QV71mezLi3WCBmGX5UaxY9ZFIyqpCSRGbdmSVfapackkm` +
		`Baqehme/N9p7+EKNTQVea71CNx9KtyKEJucKXErNPJ1iPXlZ6JZdaK6yI9Sy9Ujvjcbm` +
		`D5ULm+rJZK+uSblVjCd15klUeFA3Xpr3tfi6vTSN5nzANcG3yB4+sM8F6G45k/x9YsFV` +
		`wOb0ofzhl9kZyyX6j9v4t6/dGnr/f6zvz39Qpb7GfcZYhF7Mtzxl5K1JhMDweOEPvHLq` +
		`e17UOHbsiHD2+95UTDKpL59ffj2efn73A6Szuch3V+m9qj+n+2z+WV3u1V/vn7SfPEju` +
		`vAAwAAA==`
	_, err := Decode(input, "base64", "tar/gzip")
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}
