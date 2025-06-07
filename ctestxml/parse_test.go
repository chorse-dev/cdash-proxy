// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestXMLToJSONConversion(t *testing.T) {
	time.Local = time.UTC
	dir := "./testdata"

	xmlFiles, err := filepath.Glob(filepath.Join(dir, "*.xml"))
	if err != nil {
		t.Fatalf("Failed to glob XML files: %v", err)
	}

	for _, xmlFile := range xmlFiles {
		jsonFile := strings.TrimSuffix(xmlFile, ".xml") + ".json"

		actual, err := readXML(xmlFile)
		if err != nil {
			t.Errorf("Failed to read XML file %s: %v\n", xmlFile, err)
			continue
		}

		expected, err := readJSON(jsonFile)
		if err != nil {
			t.Errorf("Failed to read JSON file %s: %v\n", jsonFile, err)
			continue
		}

		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Errorf("Mismatch between expected and actual JSON for file %s (-expected +actual):\n%s", xmlFile, diff)
		}
	}
}

func readXML(filePath string) (any, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	job, err := Parse(file, "Example")
	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Create(strings.TrimSuffix(filePath, ".xml") + ".json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&job); err != nil {
		return nil, err
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	var data any
	err = json.Unmarshal(jobJSON, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readJSON(filePath string) (any, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var job any
	if err := json.NewDecoder(file).Decode(&job); err != nil {
		return nil, err
	}

	return job, nil
}
