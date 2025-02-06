// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/purpleKarrot/cdash-proxy/model"
)

var cfgDiagRegex = regexp.MustCompile(`CMake (Deprecation Warning|Error|Warning \(dev\)|Warning) at ([^:]+):([0-9]+) \((.*)\):`)
var cfgNameRegex = regexp.MustCompile(`"-DCMAKE_BUILD_TYPE:STRING=([^"]+)"`)

func parseConfigure(cfg *Configure, generator string) TimedCommands {
	ret := TimedCommands{
		StartTime: time.Unix(cfg.StartTime, 0),
		EndTime:   time.Unix(cfg.EndTime, 0),
	}
	ret.Commands = append(ret.Commands, model.Command{
		Name:        configName(cfg.Command),
		Type:        model.Configure,
		Status:      configureStatus(cfg.Status),
		CommandLine: cfg.Command,
		Duration:    ret.EndTime.Sub(ret.StartTime),
		Output:      cfg.Log,
		Diagnostics: splitCMakeOutput(cfg.Log),
		Attributes:  map[string]string{"Generator": generator},
		Measurements: map[string]float64{
			"Execution Time": ret.EndTime.Sub(ret.StartTime).Seconds(),
		},
	})
	return ret
}

// TODO: Add to CTest a way to specify and report the configuration to CDash.
// Example: Add a CONFIGURATION argument to cmake_configure, that results in
// "-DCMAKE_BUILD_TYPE:STRING=Debug" being passed to cmake and also adds
// <Configure><Configuration>Debug</Configuration></Configure> to the XML
// that can be used here instead of parsing the command line.
func configName(cmd string) string {
	if match := cfgNameRegex.FindStringSubmatch(cmd); match != nil {
		return fmt.Sprintf("Configure (%s)", match[1])
	}
	return "Configure"
}

func configureStatus(s int) model.CommandStatus {
	if s == 0 {
		return model.Passed
	}
	return model.Failed
}

func splitCMakeOutput(log string) []model.Diagnostic {
	var diags []model.Diagnostic
	diag := model.Diagnostic{} // TODO: set to nil

	for _, line := range strings.Split(log, "\n") {
		if len(line) == 0 {
			diag.Message += "\n"
			continue
		} else if strings.HasPrefix(line, "  ") {
			diag.Message += line[2:] + "\n"
			continue
		} else if len(diag.Message) != 0 && len(diag.FilePath) != 0 {
			diag.Message = strings.TrimRight(diag.Message, "\n")
			diags = append(diags, diag)
			diag = model.Diagnostic{} // TODO: set to nil
		}

		if match := cfgDiagRegex.FindStringSubmatch(line); match != nil {
			linenr, _ := strconv.Atoi(match[3])
			diag = model.Diagnostic{
				FilePath: match[2],
				Line:     linenr,
				Column:   -1,
				Type:     parseDiagnosticType(match[1]),
				Option:   match[4],
			}
		}
	}
	return diags
}
