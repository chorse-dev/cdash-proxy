// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/purpleKarrot/cdash-proxy/algorithm"
	"github.com/purpleKarrot/cdash-proxy/ctestxml/buildparser"
	"github.com/purpleKarrot/cdash-proxy/model"
)

var buildConfigRegex = regexp.MustCompile(`--config "([^"]+)"`)

func parseBuild(build *Build) TimedCommands {
	ret := TimedCommands{
		StartTime: time.Unix(build.StartTime, 0),
		EndTime:   time.Unix(build.EndTime, 0),
	}
	configuration := buildConfig(build.Command)
	ret.Commands = append(ret.Commands, model.Command{
		Name:        buildName("Build", configuration),
		Type:        model.Build,
		Status:      model.Passed, // TODO: model.Failed if there are 'Error's
		CommandLine: build.Command,
		Output:      combineOutput(build.Diagnostics),
		Diagnostics: mapDiagnostics(build.Diagnostics),
		Measurements: map[string]float64{
			"Execution Time": ret.EndTime.Sub(ret.StartTime).Seconds(),
		},
	})
	for _, failure := range build.Failures {
		ret.Commands = append(ret.Commands, model.Command{
			Name:        buildName(failure.Name(), configuration),
			Type:        model.Build,
			Status:      failure.Status(),
			CommandLine: strings.Join(failure.Argv, " "),
			Output:      buildparser.CleanOutput(failure.StdErr),
			Labels:      failure.Labels,
			Diagnostics: buildparser.ParseOutput(failure.SourceFile, failure.StdErr),
			Attributes:  failureVariables(&failure),
		})
	}
	return ret
}

// TODO: Add to CTest a way to report the configuration to CDash.
func buildConfig(cmd string) string {
	if match := buildConfigRegex.FindStringSubmatch(cmd); match != nil {
		return match[1]
	}
	return ""
}

func buildName(name, cfg string) string {
	if cfg != "" {
		return fmt.Sprintf("%s (%s)", name, cfg)
	}
	return name
}

func (f *Failure) Name() string {
	return fmt.Sprintf("%s while building %s %s '%s' in target %s",
		f.Type, f.Language, f.OutputType, f.OutputFile, f.Target)
}

func (f *Failure) Status() model.CommandStatus {
	if f.Type == "Error" {
		return model.Failed
	}
	return model.Passed
}

func failureVariables(failure *Failure) map[string]string {
	return map[string]string{
		"SourceFile":       failure.SourceFile,
		"WorkingDirectory": failure.WorkingDirectory,
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
		return model.Diagnostic{
			FilePath: e.SourceFile,
			Line:     e.SourceLine,
			Column:   -1,
			Type:     parseDiagnosticType(e.XMLName.Local),
			Message:  e.Text,
			Option:   "",
		}
	})
}

func parseDiagnosticType(s string) model.DiagnosticType {
	if s == "Error" {
		return model.Error
	}
	return model.Warning
}
