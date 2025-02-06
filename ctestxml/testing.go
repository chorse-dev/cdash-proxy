// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"mime"
	"strconv"
	"strings"
	"time"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/model"
)

func parseTesting(tst *Testing, sub []Subproject) TimedCommands {
	ret := TimedCommands{
		StartTime: time.Unix(tst.StartTime, 0),
		EndTime:   time.Unix(tst.EndTime, 0),
	}
	for _, t := range tst.Tests {
		cmd := model.Command{
			Name:         t.Name,
			Type:         model.Test,
			Status:       testStatus(t.Status, t.Measurements),
			CommandLine:  t.Command,
			Output:       t.Output.string,
			Labels:       t.Labels,
			Diagnostics:  parseTestOutput(t.Output.string),
			Attributes:   map[string]string{"WorkingDirectory": t.Path},
			Measurements: map[string]float64{},
		}
		for _, m := range t.Measurements {
			transformMeasurement(m, &cmd)
		}
		if p := getSubproject(sub, t.Labels); len(p) != 0 {
			cmd.Attributes["Subproject"] = p
		}
		ret.Commands = append(ret.Commands, cmd)
	}
	return ret
}

// CTest reports NOT_RUN as "notrun", COMPLETED as "passed", and all others as "passed".
// The actual status is reported the named measurement called "Exit Code" (search for
// "GetTestStatus" in the CTest source code). Note that instead of OTHER_FAULT, it
// returns the ExceptionStatus variable, which is set by the "GetExitExceptionString"
// function and may contain to many variants to store in an enum.

// enum
// { // Program statuses
//   NOT_RUN = 0,
//   TIMEOUT,
//   SEGFAULT,
//   ILLEGAL,
//   INTERRUPT,
//   NUMERICAL,
//   OTHER_FAULT,
//   FAILED,
//   BAD_COMMAND,
//   COMPLETED
// };

func testStatus(s string, measurements []Measurement) model.CommandStatus {
	switch s {
	case "passed":
		return model.Passed

	case "notrun":
		if algorithm.AnyOf(measurements, isDisabled) {
			return model.Disabled
		}
		return model.NotRun

	default:
		return model.Failed
	}
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
		cmd.Duration, _ = time.ParseDuration(string(m.Value) + "s")
		// no return. We still want the measurement!
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

func isDisabled(m Measurement) bool {
	return m.Name == "Completion Status" && string(m.Value) == "Disabled"
}
