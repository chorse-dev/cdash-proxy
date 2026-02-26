// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"bytes"
	"path/filepath"
	"strings"
	"time"

	"github.com/chorse-dev/cdash-proxy/algorithm"
	"github.com/chorse-dev/cdash-proxy/ctestxml/buildparser"
	"github.com/chorse-dev/cdash-proxy/model"
)

func parseBuild(build *Build) TimedCommands {
	startTime := time.Unix(build.StartBuildTime, 0)
	endTime := time.Unix(build.EndBuildTime, 0)
	var cmds []model.Command

	for _, command := range build.Commands.Commands {
		startTime = time.UnixMilli(command.TimeStart)
		endTime = time.UnixMilli(command.TimeStart + command.Duration)

		cmd := model.Command{
			Role:             command.Role(),
			CommandLine:      cmdFromString(command.CommandLine),
			WorkingDirectory: command.WorkingDirectory,
			StartTime:        algorithm.NewPointer(startTime),
			Duration:         command.Duration,
			Measurements:     map[string]float64{},
		}

		transformMeasurements(command.Measurements, &cmd)

		cmds = append(cmds, cmd)
	}

	if len(build.Commands.Commands) == 0 {
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
			source := stripSourcePath(command.Source, build.BinaryDirectory, build.SourceDirectory)
			cmd := model.Command{
				Role:             command.Role(),
				Result:           command.Result,
				CommandLine:      cmdFromString(command.CommandLine),
				WorkingDirectory: command.WorkingDirectory,
				StartTime:        algorithm.NewPointer(time.UnixMilli(command.TimeStart)),
				Duration:         command.Duration,
				Config:           command.Config,
				Language:         command.Language,
				Source:           source,
				Target:           target.Name,
				TargetType:       target.Type,
				TargetLabels:     target.Labels,
				Attributes:       map[string]string{},
				Measurements:     map[string]float64{},
			}
			transformMeasurements(command.Measurements, &cmd)
			cmds = append(cmds, cmd)
		}
	}

	lookupTable := make(map[string]*model.Command)
	for i := range cmds {
		lookupTable[cmds[i].CommandLine] = &cmds[i]
	}

	for _, failure := range build.Failures {
		commandLine := failure.CommandLine()
		stdout := failure.CleanStdOut()
		stderr := failure.CleanStdErr()
		diagnostics := failure.Diagnostics()

		if cmd, found := lookupTable[commandLine]; found {
			cmd.StdOut = stdout
			cmd.StdErr = stderr
			cmd.Diagnostics = diagnostics
			continue
		}

		cmds = append(cmds, model.Command{
			CommandLine:      commandLine,
			WorkingDirectory: failure.WorkingDirectory,
			Result:           failure.ExitCondition,
			Role:             "compile",
			Target:           failure.Target,
			TargetType:       failure.OutputType,
			TargetLabels:     failure.Labels,
			Source:           failure.SourceFile,
			Language:         failure.Language,
			StdOut:           stdout,
			StdErr:           stderr,
			Diagnostics:      diagnostics,
			Attributes:       map[string]string{},
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
		diag := model.Diagnostic{
			FilePath: e.SourceFile,
			Line:     e.SourceLine,
			Column:   -1,
			Type:     e.XMLName.Local,
			Message:  e.Text,
		}
		if opt := buildparser.ParseDiagnostic(e.Text); opt != nil {
			diag.Column = opt.Column
			diag.Message = opt.Message
			diag.Option = opt.Option
		}
		return diag
	})
}

func stripSourcePath(file, build, root string) string {
	if file == "" {
		return file
	}

	if build != "" {
		if rel, ok := relIfUnderDir(build, file); ok {
			return "<build>/" + rel
		}
	}

	if root != "" {
		if rel, ok := relIfUnderDir(root, file); ok {
			return rel
		}
	}

	return file
}

func relIfUnderDir(baseDir, target string) (string, bool) {
	rel, err := filepath.Rel(baseDir, target)
	if err != nil {
		return "", false
	}

	rel = filepath.Clean(rel)
	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}

	return filepath.ToSlash(rel), true
}
