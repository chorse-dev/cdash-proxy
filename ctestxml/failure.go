// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chorse-dev/cdash-proxy/algorithm"
	"github.com/chorse-dev/cdash-proxy/ctestxml/buildparser"
	"github.com/chorse-dev/cdash-proxy/model"
)

var failure_log_replacer = strings.NewReplacer(
	"[CTest: warning suppressed] ", "",
	"[CTest: warning matched] ", "",
)

func (f *Failure) CleanStdOut() string {
	return failure_log_replacer.Replace(f.StdOut)
}

func (f *Failure) CleanStdErr() string {
	return failure_log_replacer.Replace(f.StdErr)
}

func (f *Failure) CommandLine() string {
	args := algorithm.Map(f.Argv, func(s string) string {
		if strings.Contains(s, " ") {
			return strconv.Quote(s)
		}
		return s
	})

	return strings.Join(args, " ")
}

func (f *Failure) Diagnostics() []model.Diagnostic {
	var diags []model.Diagnostic
	for _, line := range strings.Split(f.CleanStdErr(), "\n") {
		if opt := buildparser.ParseDiagnostic(line); opt != nil {
			diags = append(diags, *opt)
		}
	}

	if len(diags) == 0 && f.ExitCondition != 0 {
		diags = append(diags, model.Diagnostic{
			Line:    -1,
			Column:  -1,
			Type:    "Error",
			Message: fmt.Sprintf("Command finished with exit code %d", f.ExitCondition),
		})
	}

	// 1. Loop over all diags, find an elem where elem.FilePath ends with file.
	elem := algorithm.FindIf(diags, func(d model.Diagnostic) bool {
		return strings.HasSuffix(d.FilePath, f.SourceFile)
	})
	if elem == nil {
		return diags
	}

	// 2. calculate the prefix
	prefix := elem.FilePath[:len(elem.FilePath)-len(f.SourceFile)]

	// 3. Loop over all diags, remove the prefix from all FilePaths.
	for idx, diag := range diags {
		diags[idx].FilePath = strings.TrimPrefix(diag.FilePath, prefix)
	}

	return diags
}
