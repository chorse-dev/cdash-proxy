// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package memcheck

import (
	"regexp"
	"strings"

	"github.com/purpleKarrot/cdash-proxy/model"
)

var regex = regexp.MustCompile(`^<b>([A-Z]{3})</b>`)

func parseValgrind(log string) []model.Diagnostic {
	result := []model.Diagnostic{}
	for _, line := range strings.Split(log, "\n") {
		if match := regex.FindStringSubmatch(line); match != nil {
			result = append(result, model.Diagnostic{
				FilePath: ".",
				Line:     0,
				Column:   0,
				Type:     model.Warning, // TODO: model.Defect,
				Message:  line,
				Option:   match[1],
			})
		}
	}
	return result
}
