// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"regexp"
	"strings"
)

func cmdFromString(input string) string {
	re := regexp.MustCompile(`"((?:[^"\\]|\\.)*)"`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		s := match[1 : len(match)-1]
		if strings.Contains(s, " ") {
			return match
		}
		return s
	})
}
