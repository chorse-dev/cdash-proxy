// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/chorse-dev/cdash-proxy/model"
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
		// https://gitlab.kitware.com/cmake/cmake/-/merge_requests/9860
		site.ModelName = fmt.Sprintf("Some %s CPU", site.VendorID)
	} else {
		// https://gitlab.kitware.com/utils/kwsys/-/merge_requests/339
		site.ModelName = strings.TrimSpace(site.ModelName)
	}

	job := &model.Job{
		JobID:      GenerateJobID(project, site.Name, site.BuildStamp, site.BuildName),
		Project:    project,
		BuildName:  site.BuildName,
		BuildGroup: extractGroupFromBuildstamp(site.BuildStamp),
		ChangeID:   site.ChangeID,
	}

	if site.VendorString != "" {
		job.Host = &model.Host{
			Site:           site.Name,
			Name:           site.Hostname,
			PhysicalMemory: site.TotalPhysicalMemory,
			VirtualMemory:  site.TotalVirtualMemory,
		}
		job.Host.CPU = model.CPU{
			Vendor:        site.VendorString,
			VendorID:      site.VendorID,
			FamilyID:      site.FamilyID,
			ModelID:       site.ModelID,
			ModelName:     site.ModelName,
			LogicalCores:  site.NumberOfLogicalCPU,
			PhysicalCores: site.NumberOfPhysicalCPU,
			CacheSize:     site.ProcessorCacheSize,
		}
		job.Host.OS = model.OS{
			Name:     site.OSName,
			Release:  site.OSRelease,
			Version:  site.OSVersion,
			Platform: site.OSPlatform,
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

func extractGroupFromBuildstamp(buildstamp string) string {
	if parts := strings.SplitN(buildstamp, "-", 3); len(parts) == 3 {
		return parts[2]
	}
	return ""
}
