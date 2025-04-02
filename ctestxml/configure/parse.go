// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package configure

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/purpleKarrot/cdash-proxy/model"
)

var cfgDiagRegex = regexp.MustCompile(`CMake (Deprecation Warning|Error|Warning \(dev\)|Warning)( at ([^:]+):([0-9]+) \((.*)\))?:`)

func Parse(log string, result int) []model.Diagnostic {
	var diags []model.Diagnostic
	var diag *model.Diagnostic

	for _, line := range strings.Split(log, "\n") {
		if len(line) == 0 {
			if diag != nil {
				diag.Message += "\n"
			}
			continue
		}

		if strings.HasPrefix(line, "  ") {
			if diag != nil {
				diag.Message += line[2:] + "\n"
			}
			continue
		}

		if diag != nil {
			diag.Message = strings.TrimRight(diag.Message, "\n")
			diags = append(diags, *diag)
			diag = nil
		}

		if match := cfgDiagRegex.FindStringSubmatch(line); match != nil {
			diag = &model.Diagnostic{
				FilePath: "CMakeLists.txt",
				Line:     -1,
				Column:   -1,
				Type:     cmakeDiagnosticType(match[1]),
			}
			if match[2] != "" {
				linenr, _ := strconv.Atoi(match[4])
				diag.Line = linenr
				diag.FilePath = match[3]
				diag.Option = match[5]
			}
		}
	}

	if diag != nil {
		diag.Message = strings.TrimRight(diag.Message, "\n")
		diags = append(diags, *diag)
		diag = nil
	}

	if len(diags) == 0 && result != 0 {
		diags = append(diags, model.Diagnostic{
			FilePath: "CMakeLists.txt",
			Line:     -1,
			Column:   -1,
			Type:     "Error",
			Message:  fmt.Sprintf("Command finished with exit code %d", result),
		})
	}

	return diags
}

func cmakeDiagnosticType(s string) string {
	if s == "Error" {
		return "Error"
	}
	return "Warning"
}
