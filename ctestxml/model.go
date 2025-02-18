// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/base64"
	"encoding/xml"
	"strings"

	"github.com/purpleKarrot/cdash-proxy/ctestxml/util"
)

type Done struct {
	BuildID string `xml:"buildId"`
	Time    int64  `xml:"time"`
}

type Update struct {
	Mode       string       `xml:"mode,attr"`
	Generator  string       `xml:"Generator,attr"`
	Site       string       `xml:"Site"`
	BuildName  string       `xml:"BuildName"`
	BuildStamp string       `xml:"BuildStamp"`
	StartTime  int64        `xml:"StartTime"`
	EndTime    int64        `xml:"EndTime"`
	Command    string       `xml:"UpdateCommand"`
	Type       string       `xml:"UpdateType"`
	Revision   string       `xml:"Revision"`
	Status     string       `xml:"UpdateReturnStatus"`
	Files      []UpdateFile `xml:"Directory>Updated"`
}

type UpdateFile struct {
	FullName       string `xml:"FullName"`
	Author         string `xml:"Author"`
	Email          string `xml:"Email"`
	CheckinDate    string `xml:"CheckinDate"`
	Committer      string `xml:"Committer"`
	CommitterEmail string `xml:"CommitterEmail"`
	CommitDate     string `xml:"CommitDate"`
	Log            string `xml:"Log"`
	Revision       string `xml:"Revision"`
}

type Site struct {
	ChangeID                string           `xml:"ChangeId,attr"`
	BuildName               string           `xml:"BuildName,attr"`
	BuildStamp              string           `xml:"BuildStamp,attr"`
	Name                    string           `xml:"Name,attr"`
	Generator               string           `xml:"Generator,attr"`
	OSName                  string           `xml:"OSName,attr"`
	Hostname                string           `xml:"Hostname,attr"`
	OSRelease               string           `xml:"OSRelease,attr"`
	OSVersion               string           `xml:"OSVersion,attr"`
	OSPlatform              string           `xml:"OSPlatform,attr"`
	VendorString            string           `xml:"VendorString,attr"`
	VendorID                string           `xml:"VendorID,attr"`
	FamilyID                int              `xml:"FamilyID,attr"`
	ModelID                 int              `xml:"ModelID,attr"`
	ModelName               string           `xml:"ModelName,attr"`
	ProcessorCacheSize      int              `xml:"ProcessorCacheSize,attr"`
	NumberOfLogicalCPU      int              `xml:"NumberOfLogicalCPU,attr"`
	NumberOfPhysicalCPU     int              `xml:"NumberOfPhysicalCPU,attr"`
	TotalVirtualMemory      int              `xml:"TotalVirtualMemory,attr"`
	TotalPhysicalMemory     int              `xml:"TotalPhysicalMemory,attr"`
	ProcessorClockFrequency float32          `xml:"ProcessorClockFrequency,attr"`
	Subprojects             []Subproject     `xml:"Subproject"`
	Configure               *Configure       `xml:"Configure"`
	Build                   *Build           `xml:"Build"`
	Testing                 *Testing         `xml:"Testing"`
	Coverage                *Coverage        `xml:"Coverage"`
	CoverageLog             *CoverageLog     `xml:"CoverageLog"`
	DynamicAnalysis         *DynamicAnalysis `xml:"DynamicAnalysis"`
	Notes                   []Note           `xml:"Notes>Note"`
	Uploads                 []Upload         `xml:"Upload>File"`
}

type Subproject struct {
	Label string `xml:"Label"`
	Name  string `xml:"name,attr"`
}

type Configure struct {
	StartConfigureTime int64    `xml:"StartConfigureTime"`
	EndConfigureTime   int64    `xml:"EndConfigureTime"`
	ConfigureCommand   string   `xml:"ConfigureCommand"`
	Log                string   `xml:"Log"`
	ConfigureStatus    int      `xml:"ConfigureStatus"`
	Commands           Commands `xml:"Commands"`
}

type Build struct {
	StartBuildTime int64        `xml:"StartBuildTime"`
	EndBuildTime   int64        `xml:"EndBuildTime"`
	BuildCommand   string       `xml:"BuildCommand"`
	Diagnostics    []Diagnostic `xml:",any"`
	Failures       []Failure    `xml:"Failure"`
	Targets        []Target     `xml:"Targets>Target"`
	Commands       Commands     `xml:"Commands"`

	// Make sure they do not get catched by ",any"
	StartDateTime  struct{} `xml:"StartDateTime"`
	EndDateTime    struct{} `xml:"EndDateTime"`
	Log            struct{} `xml:"Log"`
	ElapsedMinutes struct{} `xml:"ElapsedMinutes"`
}

type Target struct {
	Name     string   `xml:"name,attr"`
	Type     string   `xml:"type,attr"`
	Labels   []string `xml:"Labels>Label"`
	Commands Commands `xml:"Commands"`
}

type Commands struct {
	Commands []Command `xml:",any"`
}

type Command struct {
	XMLName      xml.Name
	Version      int           `xml:"version,attr"`
	Command      string        `xml:"command,attr"`
	Result       int           `xml:"result,attr"`
	Target       string        `xml:"target,attr"`
	TargetType   string        `xml:"targetType,attr"`
	TimeStart    int64         `xml:"timeStart,attr"`
	Duration     int64         `xml:"duration,attr"`
	Source       string        `xml:"source,attr"`
	Language     string        `xml:"language,attr"`
	Config       string        `xml:"config,attr"`
	Measurements []Measurement `xml:"NamedMeasurement"`
}

func (c Command) Role() string {
	s := c.XMLName.Local
	return strings.ToLower(string(s[0])) + s[1:]
}

type Diagnostic struct {
	XMLName     xml.Name
	Line        int    `xml:"BuildLogLine"`
	Text        string `xml:"Text"`
	SourceFile  string `xml:"SourceFile"`
	SourceLine  int    `xml:"SourceLineNumber"`
	PreContext  string `xml:"PreContext"`
	PostContext string `xml:"PostContext"`
}

type Failure struct {
	Type             string   `xml:"type,attr"`
	Target           string   `xml:"Action>TargetName"`
	Language         string   `xml:"Action>Language"`
	SourceFile       string   `xml:"Action>SourceFile"`
	OutputFile       string   `xml:"Action>OutputFile"`
	OutputType       string   `xml:"Action>OutputType"`
	WorkingDirectory string   `xml:"Command>WorkingDirectory"`
	Argv             []string `xml:"Command>Argument"`
	StdOut           string   `xml:"Result>StdOut"`
	StdErr           string   `xml:"Result>StdErr"`
	ExitCondition    int      `xml:"Result>ExitCondition"`
	Labels           []string `xml:"Labels>Label"`
}

type Testing struct {
	StartTime int64  `xml:"StartTestTime"`
	EndTime   int64  `xml:"EndTestTime"`
	Tests     []Test `xml:"Test"`
}

type Test struct {
	Name         string        `xml:"Name"`
	Path         string        `xml:"Path"`
	FullName     string        `xml:"FullName"`
	Command      string        `xml:"FullCommandLine"`
	Status       string        `xml:"Status,attr"`
	Output       Output        `xml:"Results>Measurement>Value"`
	Measurements []Measurement `xml:"Results>NamedMeasurement"`
	Labels       []string      `xml:"Labels>Label"`
}

type Output struct {
	string
}

func (o *Output) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Encoding    string `xml:"encoding,attr"`
		Compression string `xml:"compression,attr"`
		Content     string `xml:",chardata"`
	}

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	str, err := util.Decode(v.Content, v.Encoding, v.Compression)
	if err != nil {
		return err
	}

	o.string = string(str)
	return nil
}

type Measurement struct {
	Name     string `xml:"name,attr"`
	Filename string `xml:"filename,attr"`
	Type     string `xml:"type,attr"`
	Value    []byte `xml:"Value"`
}

func (o *Measurement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Name        string `xml:"name,attr"`
		Filename    string `xml:"filename,attr"`
		Type        string `xml:"type,attr"`
		Encoding    string `xml:"encoding,attr"`
		Compression string `xml:"compression,attr"`
		Value       string `xml:"Value"`
	}

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	value, err := util.Decode(v.Value, v.Encoding, v.Compression)
	if err != nil {
		return err
	}

	o.Name = v.Name
	o.Filename = v.Filename
	o.Type = v.Type
	o.Value = value
	return nil
}

type Coverage struct {
	StartTime int64          `xml:"StartTime"`
	EndTime   int64          `xml:"EndTime"`
	Files     []CoverageFile `xml:"File"`
}

type CoverageFile struct {
	// Name           string   `xml:"Name,attr"`
	Path              string   `xml:"FullPath,attr"`
	LinesTested       *int     `xml:"LOCTested,omitempty"`
	LinesUntested     *int     `xml:"LOCUnTested,omitempty"`
	BranchesTested    *int     `xml:"BranchesTested,omitempty"`
	BranchesUntested  *int     `xml:"BranchesUnTested,omitempty"`
	FunctionsTested   *int     `xml:"FunctionsTested,omitempty"`
	FunctionsUntested *int     `xml:"FunctionsUnTested,omitempty"`
	Labels            []string `xml:"Labels>Label"`
}

type CoverageLog struct {
	StartTime int64             `xml:"StartTime"`
	EndTime   int64             `xml:"EndTime"`
	Files     []CoverageLogFile `xml:"File"`
}

type CoverageLogFile struct {
	Path  string            `xml:"FullPath,attr"`
	Lines []CoverageLogLine `xml:"Report>Line"`
}

type CoverageLogLine struct {
	Number int    `xml:"Number,attr"`
	Count  int    `xml:"Count,attr"`
	Text   string `xml:",innerxml"`
}

type DynamicAnalysis struct {
	StartTime int64                 `xml:"StartTestTime"`
	EndTime   int64                 `xml:"EndTestTime"`
	Checker   string                `xml:"Checker,attr"`
	Tests     []DynamicAnalysisTest `xml:"Test"`
}

type DynamicAnalysisTest struct {
	Status      string                  `xml:"Status,attr"`
	Name        string                  `xml:"Name"`
	Path        string                  `xml:"Path"`
	FullName    string                  `xml:"FullName"`
	CommandLine string                  `xml:"FullCommandLine"`
	Defects     []DynamicAnalysisDefect `xml:"Results>Defect"`
	Log         Output                  `xml:"Log"`
}

type DynamicAnalysisDefect struct {
	Type  string `xml:"type,attr"`
	Count int    `xml:",innerxml"`
}

type Note struct {
	Name string `xml:"Name,attr"`
	Text string `xml:"Text"`
}

type Upload struct {
	Name    string `xml:"filename,attr"`
	Content []byte `xml:"Content"`
}

func (o *Upload) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Name    string `xml:"filename,attr"`
		Content string `xml:"Content"`
	}

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	// we need 1, 2, 3, or 0 padding bytes. CTest writes 1, 2, 3, or 4
	value, err := base64.StdEncoding.DecodeString(strings.TrimSuffix(v.Content, "===="))
	if err != nil {
		return err
	}

	o.Name = v.Name
	o.Content = value
	return nil
}

type Response struct {
	XMLName xml.Name `xml:"cdash"`
	Status  string   `xml:"status"` // OK, ERROR
	Message string   `xml:"message,omitempty"`
	BuildID string   `xml:"buildId,omitempty"`
}
