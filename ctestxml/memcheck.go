// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"time"

	"github.com/chorse-dev/cdash-proxy/algorithm"
	"github.com/chorse-dev/cdash-proxy/ctestxml/memcheck"
	"github.com/chorse-dev/cdash-proxy/model"
)

func parseDynamicAnalysis(da *DynamicAnalysis) TimedCommands {
	return TimedCommands{
		StartTime: time.Unix(da.StartTime, 0),
		EndTime:   time.Unix(da.EndTime, 0),
		Commands:  transformTestsDA(da),
	}
}

func transformTestsDA(da *DynamicAnalysis) []model.Command {
	return algorithm.Map(da.Tests, func(t DynamicAnalysisTest) model.Command {
		return model.Command{
			TestName:     t.Name,
			Role:         "test",
			TestStatus:   t.Status,
			CommandLine:  t.CommandLine,
			StdOut:       t.Log.string,
			Diagnostics:  memcheck.Parse(da.Checker, t.Log.string),
			Attributes:   map[string]string{"DA Checker": da.Checker},
			Measurements: memcheckParseDefects(t.Defects)}
	})
}

func memcheckParseDefects(defects []DynamicAnalysisDefect) map[string]float64 {
	meas := map[string]float64{}
	for _, d := range defects {
		meas[d.Type] = float64(d.Count)
	}
	return meas
}
