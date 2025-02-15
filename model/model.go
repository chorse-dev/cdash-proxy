// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package model

import "time"

type Job struct {
	JobID              string         `json:"job_id"`
	Project            string         `json:"project,omitempty"`
	BuildName          string         `json:"build_name,omitempty"`
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
	Done               bool           `json:"done,omitempty"`
}

type Site struct {
	Name              string `json:"name"`
	Hostname          string `json:"hostname"`
	CPUVendor         string `json:"cpu_vendor"`
	CPUVendorID       string `json:"cpu_vendor_id"`
	CPUFamilyID       int    `json:"cpu_family_id"`
	CPUModelID        int    `json:"cpu_model_id"`
	CPUModelName      string `json:"cpu_model_name"`
	CPULogicalCores   int    `json:"cpu_logical_cores"`
	CPUPhysicalCores  int    `json:"cpu_physical_cores"`
	CPUCacheSize      int    `json:"cpu_cache_size"`
	CPUClockFrequency int    `json:"cpu_clock_frequency"`
	OSName            string `json:"os_name"`
	OSRelease         string `json:"os_release"`
	OSVersion         string `json:"os_version"`
	OSPlatform        string `json:"os_platform"`
	PhysicalMemory    int    `json:"physical_memory"`
	VirtualMemory     int    `json:"virtual_memory"`
}

type Command struct {
	CommandLine      string             `json:"command_line"`
	Result           int                `json:"result"`
	Role             string             `json:"role"`
	Target           string             `json:"target,omitempty"`
	TargetType       string             `json:"target_type,omitempty"`
	TargetLabels     []string           `json:"target_labels,omitempty"`
	StartTime        *time.Time         `json:"start_time,omitempty"`
	Duration         int64              `json:"duration,omitempty"`
	Source           string             `json:"source,omitempty"`
	Language         string             `json:"language,omitempty"`
	TestName         string             `json:"test_name,omitempty"`
	TestStatus       string             `json:"test_status,omitempty"`
	StdOut           string             `json:"stdout,omitempty"`
	StdErr           string             `json:"stderr,omitempty"`
	WorkingDirectory string             `json:"working_directory,omitempty"`
	Diagnostics      []Diagnostic       `json:"diagnostics,omitempty"`
	AttachedFiles    []AttachedFile     `json:"attached_files,omitempty"`
	Attributes       map[string]string  `json:"attributes,omitempty"`
	Measurements     map[string]float64 `json:"measurements,omitempty"`
}

type AttachedFile struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Content  []byte `json:"content"`
}

type Coverage struct {
	FilePath          string   `json:"file_path"`
	Lines             []int    `json:"lines,omitempty"`
	LinesTested       *int     `json:"lines_tested,omitempty"`
	LinesUntested     *int     `json:"lines_untested,omitempty"`
	BranchesTested    *int     `json:"branches_tested,omitempty"`
	BranchesUntested  *int     `json:"branches_untested,omitempty"`
	FunctionsTested   *int     `json:"functions_tested,omitempty"`
	FunctionsUntested *int     `json:"functions_untested,omitempty"`
	Labels            []string `json:"labels,omitempty"`
}

type Diagnostic struct {
	FilePath string `json:"file_path"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Option   string `json:"option"`
}
