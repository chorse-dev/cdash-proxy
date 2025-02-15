// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package gcovtar

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	time.Local = time.UTC
	dir := "./testdata"

	gcovFiles, err := filepath.Glob(filepath.Join(dir, "*.tbz2"))
	if err != nil {
		t.Fatalf("Failed to glob GCOV files: %v", err)
	}

	for _, gcovFile := range gcovFiles {
		jsonFile := strings.TrimSuffix(gcovFile, ".tbz2") + ".json"

		actual, err := readGCOV(gcovFile)
		if err != nil {
			t.Errorf("Failed to read GCOV file %s: %v\n", gcovFile, err)
			continue
		}

		expected, err := readJSON(jsonFile)
		if err != nil {
			t.Errorf("Failed to read JSON file %s: %v\n", jsonFile, err)
			continue
		}

		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Errorf("Mismatch between expected and actual JSON for file %s (-expected +actual):\n%s", gcovFile, diff)
		}
	}
}

func readGCOV(filePath string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	job, err := Parse(file, "Example")
	if err != nil {
		return nil, err
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(jobJSON, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readJSON(filePath string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var job interface{}
	if err := json.NewDecoder(file).Decode(&job); err != nil {
		return nil, err
	}

	return job, nil
}
