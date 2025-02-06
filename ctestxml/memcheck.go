// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"time"

	"github.com/purpleKarrot/cdash-proxy/ctestxml/memcheck"
	"github.com/purpleKarrot/cdash-proxy/model"
)

func parseDynamicAnalysis(da *DynamicAnalysis) TimedCommands {
	var ret TimedCommands
	for _, t := range da.Tests {
		ret.Commands = append(ret.Commands, model.Command{
			Name:         t.Name,
			Type:         model.Test,
			Status:       testStatus(t.Status, nil),
			CommandLine:  t.CommandLine,
			Output:       t.Log.string,
			Diagnostics:  memcheck.Parse(da.Checker, t.Log.string),
			Attributes:   map[string]string{"DA Checker": da.Checker},
			Measurements: memcheckParseDefects(t.Defects),
		})
	}
	ret.StartTime = time.Unix(da.StartTime, 0)
	ret.EndTime = time.Unix(da.EndTime, 0)
	return ret
}

func memcheckParseDefects(defects []DynamicAnalysisDefect) map[string]float64 {
	meas := map[string]float64{}
	for _, d := range defects {
		meas[d.Type] = float64(d.Count)
	}
	return meas
}
