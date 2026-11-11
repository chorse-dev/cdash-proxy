// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/xml"
	"time"

	"github.com/chorse-dev/cdash-proxy/model"
)

// parseUpdate extracts repository revisions and change identity from Update.xml.
func parseUpdate(dec *xml.Decoder, elem *xml.StartElement, project string) (*model.Job, error) {
	var update Update
	if err := dec.DecodeElement(&update, elem); err != nil {
		return nil, err
	}

	startTime := time.Unix(update.StartTime, 0)
	endTime := time.Unix(update.EndTime, 0)

	job := &model.Job{
		JobID:           GenerateJobID(project, update.Site, update.BuildStamp, update.BuildName),
		BuildName:       update.BuildName,
		BuildGroup:      extractGroupFromBuildstamp(update.BuildStamp),
		ChangeID:        update.ChangeID,
		Revision:        update.Revision,
		PriorRevision:   update.PriorRevision,
		Generator:       update.Generator,
		Project:         project,
		StartUpdateTime: &startTime,
		EndUpdateTime:   &endTime,
	}

	return job, nil
}
