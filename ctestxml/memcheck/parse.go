// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package memcheck

import "github.com/purpleKarrot/cdash-proxy/model"

func Parse(checker, log string) []model.Diagnostic {
	switch checker {
	case "Valgrind":
		return parseValgrind(log)
	case "DrMemory":
		return parseDrMemory(log)
	case "Purify":
		return parsePurify(log)
	case "BoundsChecker":
		return parseBoundsChecker(log)
	case "CudaSanitizer":
		return parseCudaSanitizer(log)
	case "AddressSanitizer":
		return parseAddressSanitizer(log)
	case "LeakSanitizer":
		return parseLeakSanitizer(log)
	case "ThreadSanitizer":
		return parseThreadSanitizer(log)
	case "MemorySanitizer":
		return parseMemorySanitizer(log)
	case "UndefinedBehaviorSanitizer":
		return parseUBSanitizer(log)
	}
	return []model.Diagnostic{}
}
