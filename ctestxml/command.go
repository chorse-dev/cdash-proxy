// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
)

func cmdFromString(input string) string {
	re := regexp.MustCompile(`"((?:[^"\\]|\\.)*)"`)
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		s := match[1 : len(match)-1]
		if strings.Contains(s, " ") {
			return match
		}
		return s
	})

	fmt.Printf("CMD: %s\n", result)
	return result
}

func cmdFromArgv(argv []string) string {
	args := algorithm.Map(argv, func(s string) string {
		if strings.Contains(s, " ") {
			return strconv.Quote(s)
		}
		return s
	})

	fmt.Printf("CMD: %s\n", strings.Join(args, " "))
	return strings.Join(args, " ")
}
