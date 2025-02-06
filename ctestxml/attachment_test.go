// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"os"
	"testing"
)

func TestNotes(t *testing.T) {
	f, err := os.Open("testdata/Notes.xml")
	if err != nil {
		t.Fatal(err)
	}

	model, err := Parse(f, "purpleKarrot/Example")
	if err != nil {
		t.Fatal(err)
	}

	if len(model.AttachedFiles) != 1 {
		t.Fatal("Expected one attached file")
	}
}

func TestUpload(t *testing.T) {
	f, err := os.Open("testdata/Upload.xml")
	if err != nil {
		t.Fatal(err)
	}

	model, err := Parse(f, "purpleKarrot/Example")
	if err != nil {
		t.Fatal(err)
	}

	if len(model.AttachedFiles) != 1 {
		t.Fatal("Expected one attached file")
	}
}
