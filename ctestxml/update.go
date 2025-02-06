// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/xml"
	"time"

	"github.com/purpleKarrot/cdash-proxy/model"
)

// The only relevant piece of information that we parse from Update.xml is the CommitID.
// If there is a way to set CTEST_CHANGE_ID, then submitting Update.xml is not necessary.
// This should be the case for builds that are triggered through github actions.
// Updating requires Write Access to the source directory.
// We may set up a CI server that updates, and then invokes CTest with the source directory mounted as Read-Only.

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
		ChangeID:        update.Revision,
		Project:         project,
		StartUpdateTime: &startTime,
		EndUpdateTime:   &endTime,
	}

	return job, nil
}
