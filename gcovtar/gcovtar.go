// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package gcovtar

import (
	"archive/tar"
	"bufio"
	"compress/bzip2"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/purpleKarrot/cdash-proxy/model"
)

/**
 * Parse a tarball of .gcov files.
 **/
func Parse(r io.Reader, jobID string) (*model.Job, error) {
	tr := tar.NewReader(bzip2.NewReader(r))
	h := gcovTarHandler{}

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		fmt.Printf("-- Parsing file %s\n", hdr.Name)

		if filepath.Base(hdr.Name) == "data.json" {
			if err := h.ParseDataFile(tr); err != nil {
				return nil, err
			}
		} else if filepath.Base(hdr.Name) == "Labels.json" {
			if err := h.ParseLabelsFile(tr); err != nil {
				return nil, err
			}
		} else if filepath.Ext(hdr.Name) == ".gcov" {
			h.ParseGcovFile(tr)
		}
	}

	var cleanCoverace []model.Coverage
	for _, c := range h.Coverage {
		if strings.HasPrefix(c.FilePath, h.SourceDirectory) {
			c.FilePath = c.FilePath[len(h.SourceDirectory):]
		} else {
			continue
		}
		cleanCoverace = append(cleanCoverace, c)
	}
	h.Coverage = cleanCoverace

	// Insert coverage summary (removing any old results first)
	////$this->CoverageSummary->RemoveAll();
	//$this->CoverageSummary->Insert();
	//$this->CoverageSummary->ComputeDifference();

	//// If this source file isn't from the source or binary directory
	//// we shouldn't include it in our coverage report.
	//if strings.Contains(path, h.SourceDirectory) {
	//	path = strings.Replace(path, h.SourceDirectory, ".", -1)
	//} else if strings.Contains(path, h.BinaryDirectory) {
	//	path = strings.Replace(path, h.BinaryDirectory, ".", -1)
	//} else {
	//	return
	//}

	var cleanDiag []model.Diagnostic
	for _, d := range h.Command.Diagnostics {
		// if strings.HasPrefix(d.FilePath, h.BinaryDirectory) {
		// 	// TODO: mark as generated file,
		// 	// as it cannot be retrieved from github
		// 	// d.FilePath = d.FilePath[len(h.BinaryDirectory):]
		// } else
		if strings.HasPrefix(d.FilePath, h.SourceDirectory) {
			d.FilePath = d.FilePath[len(h.SourceDirectory):]
		} else {
			continue
		}
		cleanDiag = append(cleanDiag, d)
	}
	h.Command.Diagnostics = cleanDiag

	s := &model.Job{
		JobID:    jobID,
		Commands: []model.Command{h.Command},
		Coverage: h.Coverage,
	}
	return s, nil
}

type gcovTarHandler struct {
	//BuildId int64
	//CoverageSummary = new CoverageSummary();
	//CoverageSummary->BuildId = this->BuildId;
	SourceDirectory string
	BinaryDirectory string
	Coverage        []model.Coverage
	Command         model.Command
	Labels          map[string][]string
}

type jsonLabels struct {
	Target  jsonTarget   `json:"target"`
	Sources []jsonSource `json:"sources"`
}

type jsonTarget struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

type jsonSource struct {
	File   string   `json:"file"`
	Labels []string `json:"labels"`
}

type jsonData struct {
	BinaryDirectory string `json:"Binary"`
	SourceDirectory string `json:"Source"`
}

/**
 * Parse the Labels.json file.
 **/
func (h *gcovTarHandler) ParseLabelsFile(r io.Reader) error {
	var labels jsonLabels
	if err := json.NewDecoder(r).Decode(&labels); err != nil {
		return err
	}
	// for _, src := range labels.Sources {
	// 	h.Labels[src.File] = append(labels.Target.Labels, src.Labels...)
	// }
	return nil
}

/**
 * Parse the data.json file.
 **/
func (h *gcovTarHandler) ParseDataFile(r io.Reader) error {
	var data jsonData
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return err
	}
	h.SourceDirectory = path.Clean(data.SourceDirectory) + "/"
	h.BinaryDirectory = path.Clean(data.BinaryDirectory) + "/"
	return nil
}

/**
 * Parse an individual .gcov file.
 **/
func (h *gcovTarHandler) ParseGcovFile(r io.Reader) {
	scanner := bufio.NewScanner(r)

	// Begin parsing this file.
	// The first thing we look for is the full path to this source file.
	var path string
	for scanner.Scan() {
		line := scanner.Text()
		if pos := strings.Index(line, ":Source:"); pos != -1 {
			path = line[pos+8:]
			break
		}
	}

	covFile := model.Coverage{
		FilePath:         path,
		LinesTested:      new(int),
		LinesUntested:    new(int),
		BranchesTested:   new(int),
		BranchesUntested: new(int),
	}

	coveredBranches := int(0)
	uncoveredBranches := int(0)
	throwBranches := int(0)
	fallthroughBranches := int(0)
	branchText := ""

	// The lack of rewind is intentional.
	for scanner.Scan() {
		line := scanner.Text()

		// "Ordinary" entries in a .gcov file take the following format:
		// <lineNumber>: <timesHit>: <source code at that line>
		// So we check if this line matches the format & parse the
		// data out of it if so.
		// Otherwise we read through a block of these lines that doesn't
		// follow this format.  Such lines typically contain branch or
		// function level coverage data.

		fields := strings.SplitN(line, ":", 3)
		if len(fields) > 2 {
			lineNumber, _ := strconv.Atoi(strings.Trim(fields[1], " "))
			if lineNumber == 0 {
				continue
			}

			// Don't report branch coverage for this line if we only
			// encountered (throw) and (fallthrough) branches here.
			totalBranches := coveredBranches + uncoveredBranches
			if totalBranches > 0 && totalBranches > (throwBranches+fallthroughBranches) {
				*covFile.BranchesTested += coveredBranches
				*covFile.BranchesUntested += uncoveredBranches
				h.Command.Diagnostics = append(h.Command.Diagnostics, model.Diagnostic{
					FilePath: path,
					Line:     lineNumber - 1,
					Column:   -1,
					Type:     "Warning",
					Message:  branchText,
					Option:   "Branch Coverage",
				})
			}

			coveredBranches = 0
			uncoveredBranches = 0
			throwBranches = 0
			fallthroughBranches = 0
			branchText = ""

			line_coverage := 0

			timesHit := strings.Trim(fields[0], " ")
			if timesHit == "-" {
				// This is how gcov indicates a line of unexecutable code.
				line_coverage = -1
			} else if timesHit == "#####" {
				// This is how gcov indicates an uncovered line.
				line_coverage = 0
				*covFile.LinesUntested++
			} else {
				count, _ := strconv.Atoi(timesHit)
				line_coverage = count
				*covFile.LinesTested++
			}

			covFile.Lines = append(covFile.Lines, line_coverage)

		} else if strings.HasPrefix(line, "branch") {
			branchText += line + "\n"

			// Figure out whether this branch was covered or not.
			if strings.Contains(line, "taken 0%") {
				uncoveredBranches++
			} else {
				coveredBranches++
			}

			// Also keep track of the different types of branches encountered.
			if strings.Contains(line, "(throw)") {
				throwBranches++
			} else if strings.Contains(line, "(fallthrough)") {
				fallthroughBranches++
			}
		}
	}

	h.Coverage = append(h.Coverage, covFile)
}
