// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/purpleKarrot/cdash-proxy/model"
)

var cfgDiagRegex = regexp.MustCompile(`CMake (Deprecation Warning|Error|Warning \(dev\)|Warning) at ([^:]+):([0-9]+) \((.*)\):`)

func parseConfigure(cfg *Configure, generator string) TimedCommands {
	ret := TimedCommands{
		StartTime: time.Unix(cfg.StartTime, 0),
		EndTime:   time.Unix(cfg.EndTime, 0),
	}
	ret.Commands = append(ret.Commands, model.Command{
		Role:         "configure",
		Result:       cfg.Status,
		CommandLine:  cfg.Command,
		StdOut:       cfg.Log,
		StartTime:    &ret.StartTime,
		Duration:     ret.EndTime.Sub(ret.StartTime).Milliseconds(),
		Diagnostics:  splitCMakeOutput(cfg.Log),
		Attributes:   map[string]string{"Generator": generator},
		Measurements: map[string]float64{},
	})
	return ret
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
