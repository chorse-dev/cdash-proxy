// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/ctestxml/configure"
	"github.com/purpleKarrot/cdash-proxy/model"
)

type cfgStep struct {
	role  string
	regex *regexp.Regexp
}

var cfgSteps = []cfgStep{
	{"configure", regexp.MustCompile(`-- Configuring done \(([0-9.]+)s\)\n`)},
	{"generate", regexp.MustCompile(`-- Generating done \(([0-9.]+)s\)\n`)},
}

func parseConfigure(cfg *Configure, generator string) TimedCommands {
	input := cfg.Log
	startTime := time.Unix(cfg.StartConfigureTime, 0)
	endTime := time.Unix(cfg.EndConfigureTime, 0)
	duration := time.Duration(0)

	var cmds []model.Command
	for _, step := range cfgSteps {
		cmd := model.Command{
			Role:         step.role,
			CommandLine:  cfg.ConfigureCommand,
			StartTime:    algorithm.NewPointer(startTime.Add(duration)),
			Attributes:   map[string]string{"Generator": generator},
			Measurements: map[string]float64{},
		}

		match := step.regex.FindStringSubmatchIndex(input)
		if match != nil {
			seconds, _ := strconv.ParseFloat(input[match[2]:match[3]], 64)

			cmd.Result = 0
			cmd.StdOut = input[:match[1]]
			cmd.Duration = int64(seconds * 1000.0)

			input = input[match[1]:]
			duration += time.Duration(seconds * float64(time.Second))
		} else {
			cmd.Result = cfg.ConfigureStatus
			cmd.StdOut = input
			cmd.Duration = endTime.Sub(*cmd.StartTime).Milliseconds()
		}

		cmd.Diagnostics = configure.Parse(cmd.StdOut)
		cmds = append(cmds, cmd)

		if match == nil {
			break
		}
	}

	if cfg.Commands.Commands == nil {
		goto ret
	}

	if len(cmds) != len(cfg.Commands.Commands) {
		log.Printf("Command lenght mismatch: %d != %d\n", len(cmds), len(cfg.Commands.Commands))
		goto ret
	}

	for i, in := range cfg.Commands.Commands {
		out := cmds[i]

		if out.Role != in.Role() {
			log.Printf("Role mismatch: %s != %s\n", out.Role, in.Role())
		}

		out.Result = in.Result
		out.StartTime = algorithm.NewPointer(time.UnixMilli(in.TimeStart))
		out.Duration = in.Duration
		transformMeasurements(in.Measurements, &out)
		cmds[i] = out
	}

ret:
	return TimedCommands{
		StartTime: startTime,
		EndTime:   endTime,
		Commands:  cmds,
	}
}
