// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/xml"

	"github.com/purpleKarrot/cdash-proxy/model"
)

func parseDone(dec *xml.Decoder, elem *xml.StartElement, project string) (*model.Job, error) {
	var done Done
	if err := dec.DecodeElement(&done, elem); err != nil {
		return nil, err
	}

	job := &model.Job{
		JobID:   done.BuildID,
		Project: project,
		Done:    true,
	}

	return job, nil
}
