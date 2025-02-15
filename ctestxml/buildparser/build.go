// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package buildparser

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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

var replacer = strings.NewReplacer(
	"[CTest: warning suppressed] ", "",
	"[CTest: warning matched] ", "",
)

func CleanOutput(output string) string {
	return replacer.Replace(output)
}

func ParseOutput(file, output string) []model.Diagnostic {
	if len(file) == 0 {
		// this is the case for linker errors
		return []model.Diagnostic{}
	}

	var diags []model.Diagnostic
	for _, line := range strings.Split(output, "\n") {
		kind := detectLineType(line)
		if kind == lineTypeRegular {
			continue
		}
		line = CleanOutput(line)
		diags = append(diags, ParseDiagnostic(file, kind.DiagnosticType(), line))
	}

	// TODO: We need a better way to strip the source directory!
	// 1. Loop over all diags, find an elem where elem.FilePath ends with file.
	// 2. calculate the prefix
	// 3. Loop over all diags, remove the prefix from all FilePaths.

	return diags
}

type lineType int

const (
	lineTypeRegular lineType = iota
	lineTypeError
	lineTypeWarning
	lineTypeNote
)

func (i lineType) DiagnosticType() string {
	switch i {
	default:
		return "Error"
	case lineTypeWarning:
		return "Warning"
	case lineTypeNote:
		return "Note"
	}
}

func detectLineType(line string) lineType {
	if strings.HasPrefix(line, "[CTest: warning suppressed]") {
		return lineTypeRegular
	}
	if strings.HasPrefix(line, "[CTest: warning matched]") {
		return lineTypeWarning
	}
	return lineTypeRegular
}

func ParseDiagnostic(file string, kind string, line string) model.Diagnostic {
	diag := model.Diagnostic{
		FilePath: file,
		Line:     -1,
		Column:   -1,
		Type:     kind,
		Message:  line,
		Option:   "",
	}
	for _, re := range reFileLine {
		match := re.FindStringSubmatch(line)
		if match == nil {
			continue
		}
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

		break
	}

	// TODO: We need a better way to strip the source directory!
	if strings.HasSuffix(diag.FilePath, file) {
		diag.FilePath = file
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
