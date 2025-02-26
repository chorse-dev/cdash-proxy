// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"bytes"
	"time"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/ctestxml/buildparser"
	"github.com/purpleKarrot/cdash-proxy/model"
)

func parseBuild(build *Build) TimedCommands {
	startTime := time.Unix(build.StartBuildTime, 0)
	endTime := time.Unix(build.EndBuildTime, 0)
	var cmds []model.Command

	if len(build.Commands.Commands) == 1 {
		command := build.Commands.Commands[0]
		startTime = time.UnixMilli(command.TimeStart)
		endTime = time.UnixMilli(command.TimeStart + command.Duration)

		cmd := model.Command{
			Role:         command.Role(),
			CommandLine:  cmdFromString(command.Command),
			StartTime:    algorithm.NewPointer(startTime),
			Duration:     command.Duration,
			Measurements: map[string]float64{},
		}

		transformMeasurements(command.Measurements, &cmd)

		cmds = append(cmds, cmd)
	} else {
		cmds = append(cmds, model.Command{
			Role:         "cmakeBuild",
			CommandLine:  cmdFromString(build.BuildCommand),
			StartTime:    algorithm.NewPointer(startTime),
			Duration:     endTime.Sub(startTime).Milliseconds(),
			Measurements: map[string]float64{},
		})
	}

	cmds[0].StdOut = combineOutput(build.Diagnostics)
	cmds[0].Diagnostics = mapDiagnostics(build.Diagnostics)

	for _, target := range build.Targets {
		for _, command := range target.Commands.Commands {
			cmd := model.Command{
				Role:         command.Role(),
				Result:       command.Result,
				CommandLine:  cmdFromString(command.Command),
				StartTime:    algorithm.NewPointer(time.UnixMilli(command.TimeStart)),
				Duration:     command.Duration,
				Config:       command.Config,
				Language:     command.Language,
				Source:       command.Source,
				Target:       target.Name,
				TargetType:   target.Type,
				TargetLabels: target.Labels,
				Attributes:   map[string]string{},
				Measurements: map[string]float64{},
			}
			transformMeasurements(command.Measurements, &cmd)
			cmds = append(cmds, cmd)
		}
	}

	// TODO: Merge into existing commands
	for _, failure := range build.Failures {
		cmds = append(cmds, model.Command{
			CommandLine:      cmdFromArgv(failure.Argv),
			Result:           failure.ExitCondition,
			Role:             "compile",
			Target:           failure.Target,
			TargetType:       failure.OutputType,
			TargetLabels:     failure.Labels,
			Source:           failure.SourceFile,
			Language:         failure.Language,
			StdOut:           buildparser.CleanOutput(failure.StdOut),
			StdErr:           buildparser.CleanOutput(failure.StdErr),
			Diagnostics:      buildparser.ParseOutput(failure.SourceFile, failure.StdErr),
			Attributes:       map[string]string{},
			WorkingDirectory: failure.WorkingDirectory,
			// Outputs:          failure.OutputFile,
		})
	}

	return TimedCommands{
		StartTime: startTime,
		EndTime:   endTime,
		Commands:  cmds,
	}
}

func combineOutput(messages []Diagnostic) string {
	var buffer bytes.Buffer
	for _, e := range messages {
		buffer.WriteString(e.PreContext)
		buffer.WriteString(e.Text)
		buffer.WriteString("\n")
		buffer.WriteString(e.PostContext)
	}
	return buffer.String()
}

func mapDiagnostics(messages []Diagnostic) []model.Diagnostic {
	return algorithm.Map(messages, func(e Diagnostic) model.Diagnostic {
		diag := buildparser.ParseDiagnostic(e.SourceFile, e.XMLName.Local, e.Text)
		diag.Line = e.SourceLine
		return diag
	})
}
