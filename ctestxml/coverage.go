// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"strings"
	"time"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/model"
)

func parseCoverage(cov *Coverage) (ret TimedCovarage) {
	ret.Files = algorithm.Map(cov.Files, func(f CoverageFile) model.Coverage {
		return model.Coverage{
			FilePath:          strings.TrimPrefix(f.Path, "./"),
			LinesTested:       f.LinesTested,
			LinesUntested:     f.LinesUntested,
			BranchesTested:    f.BranchesTested,
			BranchesUntested:  f.BranchesUntested,
			FunctionsTested:   f.FunctionsTested,
			FunctionsUntested: f.FunctionsUntested,
			Labels:            f.Labels,
		}
	})

	ret.StartTime = time.Unix(cov.StartTime, 0)
	ret.EndTime = time.Unix(cov.EndTime, 0)

	return
}

func parseCoverageLog(cov *CoverageLog) (ret TimedCovarage) {
	ret.Files = algorithm.Map(cov.Files, func(f CoverageLogFile) model.Coverage {
		return model.Coverage{
			FilePath: strings.TrimPrefix(f.Path, "./"),
			Lines: algorithm.Map(f.Lines, func(l CoverageLogLine) int {
				return l.Count
			}),
		}
	})

	ret.StartTime = time.Unix(cov.StartTime, 0)
	ret.EndTime = time.Unix(cov.EndTime, 0)

	return
}
