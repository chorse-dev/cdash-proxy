// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package model

import (
	"encoding/json"
	"time"
)

type Job struct {
	JobID              string         `json:"job_id"`
	Project            string         `json:"project"`
	BuildName          string         `json:"build_name"`
	ChangeID           string         `json:"change_id,omitempty"`
	Site               *Site          `json:"site,omitempty"`
	StartUpdateTime    *time.Time     `json:"start_update_time,omitempty"`
	EndUpdateTime      *time.Time     `json:"end_update_time,omitempty"`
	StartConfigureTime *time.Time     `json:"start_configure_time,omitempty"`
	EndConfigureTime   *time.Time     `json:"end_configure_time,omitempty"`
	StartBuildTime     *time.Time     `json:"start_build_time,omitempty"`
	EndBuildTime       *time.Time     `json:"end_build_time,omitempty"`
	StartTestTime      *time.Time     `json:"start_test_time,omitempty"`
	EndTestTime        *time.Time     `json:"end_test_time,omitempty"`
	StartCoverageTime  *time.Time     `json:"start_coverage_time,omitempty"`
	EndCoverageTime    *time.Time     `json:"end_coverage_time,omitempty"`
	StartMemcheckTime  *time.Time     `json:"start_memcheck_time,omitempty"`
	EndMemcheckTime    *time.Time     `json:"end_memcheck_time,omitempty"`
	Commands           []Command      `json:"commands,omitempty"`
	Coverage           []Coverage     `json:"coverage,omitempty"`
	AttachedFiles      []AttachedFile `json:"attached_files,omitempty"`
	Done               bool           `json:"done"`
}

type Site struct {
	Name           string `json:"name"`
	Hostname       string `json:"hostname"`
	CPU            CPU    `json:"cpu"`
	Kernel         Kernel `json:"kernel"`
	PhysicalMemory int    `json:"physical_memory"`
	VirtualMemory  int    `json:"virtual_memory"`
}

type CPU struct {
	Vendor         string `json:"vendor"`
	VendorID       string `json:"vendor_id"`
	FamilyID       int    `json:"family_id"`
	ModelID        int    `json:"model_id"`
	ModelName      string `json:"model_name"`
	LogicalCPUs    int    `json:"logical_cpus"`
	PhysicalCPUs   int    `json:"physical_cpus"`
	CacheSize      int    `json:"cache_size"`
	ClockFrequency int    `json:"clock_frequency"`
}

type Kernel struct {
	Name     string `json:"name"`
	Release  string `json:"release"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
}

type Command struct {
	Name          string             `json:"name"`
	Type          CommandType        `json:"type"`
	Status        CommandStatus      `json:"status"`
	Duration      time.Duration      `json:"duration"`
	CommandLine   string             `json:"command_line"`
	Output        string             `json:"output"`
	Labels        []string           `json:"labels,omitempty"`
	Diagnostics   []Diagnostic       `json:"diagnostics,omitempty"`
	AttachedFiles []AttachedFile     `json:"attached_files,omitempty"`
	Attributes    map[string]string  `json:"attributes,omitempty"`
	Measurements  map[string]float64 `json:"measurements,omitempty"`
}

type AttachedFile struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Content  []byte `json:"content"`
}

type CommandStatus int

func (t CommandStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

const (
	Passed CommandStatus = iota
	Disabled
	Failed
	NotRun
)

type CommandType int

func (t CommandType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

const (
	Configure CommandType = iota
	Build
	Test
	Valgrind
)

type Coverage struct {
	FilePath          string   `json:"file_path"`
	LineCoverage      []int    `json:"line_coverage,omitempty"`
	LinesTested       *int     `json:"lines_tested,omitempty"`
	LinesUntested     *int     `json:"lines_untested,omitempty"`
	BranchesTested    *int     `json:"branches_tested,omitempty"`
	BranchesUntested  *int     `json:"branches_untested,omitempty"`
	FunctionsTested   *int     `json:"functions_tested,omitempty"`
	FunctionsUntested *int     `json:"functions_untested,omitempty"`
	Labels            []string `json:"labels,omitempty"`
}

type Diagnostic struct {
	FilePath string         `json:"file_path"`
	Line     int            `json:"line"`
	Column   int            `json:"column"`
	Type     DiagnosticType `json:"type"`
	Message  string         `json:"message"`
	Option   string         `json:"option"`
}

type DiagnosticType int

func (t DiagnosticType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

const (
	Error DiagnosticType = iota
	Warning
	Note
	// TODO: add Defect (for static analysis)
)

//go:generate stringer -type=CommandStatus,CommandType,DiagnosticType -output=model_strings.go
