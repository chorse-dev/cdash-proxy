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

	job := &model.Job{
		JobID:     GenerateJobID(project, site.Name, site.BuildStamp, site.BuildName),
		Project:   project,
		BuildName: site.BuildName,
		ChangeID:  site.ChangeID,
	}

	if site.VendorString != "" {
		job.Site = &model.Site{
			Name:              site.Name,
			Hostname:          site.Hostname,
			CPUVendor:         site.VendorString,
			CPUVendorID:       site.VendorID,
			CPUFamilyID:       site.FamilyID,
			CPUModelID:        site.ModelID,
			CPUModelName:      site.ModelName,
			CPULogicalCores:   site.NumberOfLogicalCPU,
			CPUPhysicalCores:  site.NumberOfPhysicalCPU,
			CPUCacheSize:      site.ProcessorCacheSize,
			CPUClockFrequency: int(site.ProcessorClockFrequency),
			OSName:            site.OSName,
			OSRelease:         site.OSRelease,
			OSVersion:         site.OSVersion,
			OSPlatform:        site.OSPlatform,
			PhysicalMemory:    site.TotalPhysicalMemory,
			VirtualMemory:     site.TotalVirtualMemory,
		}
	}

	if site.Configure != nil {
		ret := parseConfigure(site.Configure, site.Generator)
		job.Commands = ret.Commands
		job.StartConfigureTime = &ret.StartTime
		job.EndConfigureTime = &ret.EndTime
	}
	if site.Build != nil {
		ret := parseBuild(site.Build)
		job.Commands = ret.Commands
		job.StartBuildTime = &ret.StartTime
		job.EndBuildTime = &ret.EndTime
	}
	if site.Testing != nil {
		ret := parseTesting(site.Testing, site.Subprojects)
		job.Commands = ret.Commands
		job.StartTestTime = &ret.StartTime
		job.EndTestTime = &ret.EndTime
	}
	if site.Coverage != nil {
		ret := parseCoverage(site.Coverage)
		job.Coverage = ret.Files
		job.StartCoverageTime = &ret.StartTime
		job.EndCoverageTime = &ret.EndTime
	}
	if site.CoverageLog != nil {
		ret := parseCoverageLog(site.CoverageLog)
		job.Coverage = ret.Files
		job.StartCoverageTime = &ret.StartTime
		job.EndCoverageTime = &ret.EndTime
	}
	if site.DynamicAnalysis != nil {
		ret := parseDynamicAnalysis(site.DynamicAnalysis)
		job.Commands = ret.Commands
		job.StartMemcheckTime = &ret.StartTime
		job.EndMemcheckTime = &ret.EndTime
	}
	if len(site.Notes) != 0 {
		job.AttachedFiles = parseNotes(site.Notes)
	}
	if len(site.Uploads) != 0 {
		job.AttachedFiles = parseUploads(site.Uploads)
	}
	return job, nil
}
