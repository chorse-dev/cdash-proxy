// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package buildparser

import (
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/model"
)

var reFileLine = algorithm.Map([]string{
	"^(?P<file>[a-zA-Z./0-9_+ ~-]+):(?P<line>[0-9]+):(?P<column>[0-9]+): (?P<type>error|warning|note): (?P<message>.*) \\[(?P<option>.*)\\]$",
	"^(?P<file>[a-zA-Z./0-9_+ ~-]+):(?P<line>[0-9]+):(?P<column>[0-9]+): (?P<type>error|warning|note): (?P<message>.*)",
	"^(?P<file>[a-zA-Z.\\:/0-9_+ ~-]+)\\((?P<line>[0-9]+)\\)",
	"^[0-9]+>(?P<file>[a-zA-Z.\\:/0-9_+ ~-]+)\\((?P<line>[0-9]+)\\)",
	"^(?P<file>[a-zA-Z./0-9_+ ~-]+)\\((?P<line>[0-9]+)\\)",
	"\"(?P<file>[a-zA-Z./0-9_+ ~-]+)\", line (?P<line>[0-9]+)",
	"File = (?P<file>[a-zA-Z./0-9_+ ~-]+), Line = (?P<line>[0-9]+)",
	"^Warning W[0-9]+ (?P<file>[a-zA-Z.\\:/0-9_+ ~-]+) (?P<line>[0-9]+):",
}, func(p string) *regexp.Regexp {
	return regexp.MustCompile(p)
})

func ParseDiagnostic(line string) *model.Diagnostic {
	for _, re := range reFileLine {
		if match := re.FindStringSubmatch(line); match != nil {
			return toDiagnostic(re, match)
		}
	}
	return nil
}

func toDiagnostic(re *regexp.Regexp, match []string) *model.Diagnostic {
	diag := &model.Diagnostic{}
	for k, name := range re.SubexpNames() {
		switch name {
		case "file":
			diag.FilePath = filepath.Clean(match[k])
		case "line":
			diag.Line, _ = strconv.Atoi(match[k])
		case "column":
			diag.Column, _ = strconv.Atoi(match[k])
		case "type":
			diag.Type = parseDiagnosticType(match[k])
		case "message":
			diag.Message = match[k]
		case "option":
			diag.Option = match[k]
		}
	}
	return diag
}

func parseDiagnosticType(s string) string {
	switch s {
	default:
		return "Error"
	case "warning":
		return "Warning"
	case "note":
		return "Note"
	}
}
