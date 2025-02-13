// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"math"
	"mime"
	"strconv"
	"strings"
	"time"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/model"
)

func parseTesting(tst *Testing, sub []Subproject) TimedCommands {
	return TimedCommands{
		StartTime: time.Unix(tst.StartTime, 0),
		EndTime:   time.Unix(tst.EndTime, 0),
		Commands:  transformTests(tst.Tests, sub),
	}
}

func transformTests(tests []Test, sub []Subproject) []model.Command {
	return algorithm.Map(tests, func(t Test) model.Command {
		cmd := model.Command{
			TestName:         t.Name,
			Role:             "Test",
			TestStatus:       t.Status,
			CommandLine:      t.Command,
			StdOut:           t.Output.string,
			TargetLabels:     t.Labels,
			WorkingDirectory: t.Path,
			Diagnostics:      parseTestOutput(t.Output.string),
			Attributes:       map[string]string{},
			Measurements:     map[string]float64{},
		}
		for _, m := range t.Measurements {
			transformMeasurement(m, &cmd)
		}
		if p := getSubproject(sub, t.Labels); len(p) != 0 {
			cmd.Attributes["Subproject"] = p
		}
		return cmd
	})
}

// https://cmake.org/cmake/help/latest/prop_test/ATTACHED_FILES.html
// https://cmake.org/cmake/help/latest/prop_test/ATTACHED_FILES_ON_FAIL.html
// https://github.com/Kitware/CMake/blob/master/Source/CTest/cmCTestTestHandler.cxx#L1451
// https://github.com/Kitware/CMake/blob/master/Source/CTest/cmCTestTestHandler.cxx#L1908
// https://github.com/Kitware/VTK/blob/master/Testing/Rendering/vtkTesting.cxx#L510

func transformMeasurement(m Measurement, cmd *model.Command) {
	if m.Name == "Command Line" {
		return // should already be in CommandLine
	}

	if m.Name == "Execution Time" {
		sec, _ := strconv.ParseFloat(string(m.Value), 64)
		cmd.Duration = int64(math.Round(sec * 1000))
		return
	}

	if m.Type == "file" {
		cmd.AttachedFiles = append(cmd.AttachedFiles, model.AttachedFile{
			Name:     m.Name,
			Filename: m.Filename,
			Type:     "application/octet-stream",
			Content:  m.Value,
		})
		return
	}

	if strings.HasPrefix(m.Type, "image/") {
		ext, err := mime.ExtensionsByType(m.Type)
		if err != nil || len(ext) == 0 {
			return
		}
		cmd.AttachedFiles = append(cmd.AttachedFiles, model.AttachedFile{
			Name:     m.Name,
			Filename: m.Name + ext[0],
			Type:     m.Type,
			Content:  m.Value,
		})
		return
	}

	if strings.HasPrefix(m.Type, "numeric/") {
		cmd.Measurements[m.Name], _ = strconv.ParseFloat(string(m.Value), 64)
		return
	}

	// case "Completion Status":
	// case "Exit Code":
	// case "Exit Value":
	// case "Fail Reason":
	// case "Pass Reason":

	cmd.Attributes[m.Name] = string(m.Value)
}

func parseTestOutput(log string) []model.Diagnostic {
	return []model.Diagnostic{}
}
