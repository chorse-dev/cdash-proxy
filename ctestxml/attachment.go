// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"net/http"
	"path/filepath"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/model"
)

func parseNotes(notes []Note) []model.AttachedFile {
	return algorithm.Map(notes, func(note Note) model.AttachedFile {
		return model.AttachedFile{
			Name:     note.Name,
			Filename: filepath.Base(note.Name),
			Type:     "text/plain",
			Content:  []byte(note.Text),
		}
	})
}

func parseUploads(uploads []Upload) []model.AttachedFile {
	return algorithm.Map(uploads, func(up Upload) model.AttachedFile {
		return model.AttachedFile{
			Name:     up.Name,
			Filename: filepath.Base(up.Name),
			Type:     http.DetectContentType(up.Content),
			Content:  up.Content,
		}
	})
}
