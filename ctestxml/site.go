// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/purpleKarrot/cdash-proxy/model"
)

type TimedCommands struct {
	Commands  []model.Command
	StartTime time.Time
	EndTime   time.Time
}

type TimedCovarage struct {
	Files     []model.Coverage
	StartTime time.Time
	EndTime   time.Time
}

func parseSite(dec *xml.Decoder, elem *xml.StartElement, project string) (*model.Job, error) {
	var site Site
	if err := dec.DecodeElement(&site, elem); err != nil {
		return nil, err
	}

	if site.ModelName == "" {
		// upstream change required!
		site.ModelName = fmt.Sprintf("%s %d %d", site.VendorID, site.FamilyID, site.ModelID)
	}

	s := &model.Job{
		JobID:     GenerateJobID(project, site.Name, site.BuildStamp, site.BuildName),
		Project:   project,
		BuildName: site.BuildName,
		ChangeID:  site.ChangeID,
		Site: &model.Site{
			Name: site.Name,
			CPU: model.CPU{
				Vendor:         site.VendorString,
				VendorID:       site.VendorID,
				FamilyID:       site.FamilyID,
				ModelID:        site.ModelID,
				ModelName:      site.ModelName,
				LogicalCPUs:    site.NumberOfLogicalCPU,
				PhysicalCPUs:   site.NumberOfPhysicalCPU,
				CacheSize:      site.ProcessorCacheSize,
				ClockFrequency: int(site.ProcessorClockFrequency),
			},
			Kernel: model.Kernel{
				Name:     site.OSName,
				Release:  site.OSRelease,
				Version:  site.OSVersion,
				Platform: site.OSPlatform,
			},
			Hostname:       site.Hostname,
			PhysicalMemory: site.TotalPhysicalMemory,
			VirtualMemory:  site.TotalVirtualMemory,
		},
	}

	if site.Configure != nil {
		ret := parseConfigure(site.Configure, site.Generator)
		s.Commands = ret.Commands
		s.StartConfigureTime = &ret.StartTime
		s.EndConfigureTime = &ret.EndTime
	}
	if site.Build != nil {
		ret := parseBuild(site.Build)
		s.Commands = ret.Commands
		s.StartBuildTime = &ret.StartTime
		s.EndBuildTime = &ret.EndTime
	}
	if site.Testing != nil {
		ret := parseTesting(site.Testing, site.Subprojects)
		s.Commands = ret.Commands
		s.StartTestTime = &ret.StartTime
		s.EndTestTime = &ret.EndTime
	}
	if site.Coverage != nil {
		ret := parseCoverage(site.Coverage)
		s.Coverage = ret.Files
		s.StartCoverageTime = &ret.StartTime
		s.EndCoverageTime = &ret.EndTime
	}
	if site.CoverageLog != nil {
		ret := parseCoverageLog(site.CoverageLog)
		s.Coverage = ret.Files
		s.StartCoverageTime = &ret.StartTime
		s.EndCoverageTime = &ret.EndTime
	}
	if site.DynamicAnalysis != nil {
		ret := parseDynamicAnalysis(site.DynamicAnalysis)
		s.Commands = ret.Commands
		s.StartMemcheckTime = &ret.StartTime
		s.EndMemcheckTime = &ret.EndTime
	}
	if len(site.Notes) != 0 {
		s.AttachedFiles = parseNotes(site.Notes)
	}
	if len(site.Uploads) != 0 {
		s.AttachedFiles = parseUploads(site.Uploads)
	}
	return s, nil
}
