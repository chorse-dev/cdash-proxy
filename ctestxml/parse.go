// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/xml"
	"errors"
	"io"

	"github.com/purpleKarrot/cdash-proxy/model"
)

func Parse(r io.Reader, project string) (*model.Job, error) {
	dec := xml.NewDecoder(r)
	se, err := startElement(dec)
	if err != nil {
		return nil, err
	}
	switch se.Name.Local {
	case "Done":
		return parseDone(dec, se, project)
	case "Site":
		return parseSite(dec, se, project)
	case "Update":
		return parseUpdate(dec, se, project)
	}
	return nil, errors.New("Unknown XML Tag " + se.Name.Local)
}

func startElement(dec *xml.Decoder) (*xml.StartElement, error) {
	for {
		t, err := dec.Token()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			return &se, nil
		case xml.EndElement:
			break
		}
	}
	return nil, errors.New("no XML Tag found")
}
