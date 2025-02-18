// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"time"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/model"
)

func parseTesting(tst *Testing, sub []Subproject) TimedCommands {
	return TimedCommands{
		StartTime: time.Unix(tst.StartTime, 0),
		EndTime:   time.Unix(tst.EndTime, 0),
		Commands:  transformTests(tst.Tests, sub),
	}
}

func transformTests(tests []Test, sub []Subproject) []model.Command {
	return algorithm.Map(tests, func(t Test) model.Command {
		cmd := model.Command{
			TestName:         t.Name,
			Role:             "test",
			TestStatus:       t.Status,
			CommandLine:      t.Command,
			StdOut:           t.Output.string,
			TargetLabels:     t.Labels,
			WorkingDirectory: t.Path,
			Diagnostics:      parseTestOutput(t.Output.string),
			Attributes:       map[string]string{},
			Measurements:     map[string]float64{},
		}
		transformMeasurements(t.Measurements, &cmd)
		if p := getSubproject(sub, t.Labels); len(p) != 0 {
			cmd.Attributes["Subproject"] = p
		}
		return cmd
	})
}

func parseTestOutput(log string) []model.Diagnostic {
	return []model.Diagnostic{}
}
