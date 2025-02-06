// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package util

import (
	"archive/tar"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"io"
	"strings"
)

func Decode(input, encoding, compression string) ([]byte, error) {
	var reader io.Reader
	reader = strings.NewReader(input)

	if encoding == "base64" {
		reader = base64.NewDecoder(base64.StdEncoding, reader)
	}

	if compression == "gzip" {
		xr, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}
		defer xr.Close()
		reader = xr
	}

	if compression == "tar/gzip" {
		xr, err := gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
		defer xr.Close()

		tr := tar.NewReader(xr)
		if _, err = tr.Next(); err != nil {
			return nil, err
		}
		reader = tr
	}

	return io.ReadAll(reader)
}
