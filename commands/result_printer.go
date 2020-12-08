package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/olekukonko/tablewriter"
	"io"
	"strings"
)

type printerConfig struct {
	printNoIssues bool
}

type resultPrinter struct {
	config printerConfig
	table  *tablewriter.Table
}
type resultLineSummary struct {
	Component   string
	HighCount   int
	MediumCount int
	LowCount    int
	FixVersions map[string]bool
	Licenses    []string
}

func (r resultLineSummary) hasIssue() bool {
	return r.LowCount > 0 || r.MediumCount > 0 || r.HighCount > 0
}

func newPrinter(writer io.Writer, config printerConfig) *resultPrinter {
	table := tablewriter.NewWriter(writer)
	table.SetRowLine(true)
	table.SetHeader([]string{"Component", "High", "Medium", "Low", "Fix version(s)", "Licenses"})
	return &resultPrinter{table: table, config: config}
}

// Consolidate scan result from Xray and prints as a table
func (r *resultPrinter) print(result ComponentSummaryResult) {
	log.Debug("Printing result for ", len(result.Artifacts), " components")
	for _, artifact := range result.Artifacts {
		lineSummary := r.createLineSummary(artifact)
		if r.config.printNoIssues || lineSummary.hasIssue() {
			r.renderLine(lineSummary)
		}
	}
}

func (r *resultPrinter) flush() {
	r.table.Render()
}

func (r *resultPrinter) renderLine(lineSummary resultLineSummary) {
	var rowColor tablewriter.Colors
	if lineSummary.LowCount == 0 && lineSummary.MediumCount == 0 && lineSummary.HighCount == 0 {
		rowColor = tablewriter.Colors{tablewriter.FgGreenColor}
	} else if lineSummary.HighCount > 0 {
		rowColor = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiMagentaColor}
	} else if lineSummary.MediumCount > 0 {
		rowColor = tablewriter.Colors{tablewriter.FgHiRedColor}
	} else if lineSummary.LowCount > 0 {
		rowColor = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	r.table.Rich([]string{
		lineSummary.Component,
		fmt.Sprintf("%d", lineSummary.HighCount),
		fmt.Sprintf("%d", lineSummary.MediumCount),
		fmt.Sprintf("%d", lineSummary.LowCount),
		versions(lineSummary.FixVersions),
		strings.Join(lineSummary.Licenses, ", ")},
		[]tablewriter.Colors{
			rowColor,
			rowColor,
			rowColor,
			rowColor,
			{tablewriter.FgHiWhiteColor},
			{tablewriter.FgHiWhiteColor},
		},
	)
}

func (r *resultPrinter) createLineSummary(artifact ComponentArtifact) resultLineSummary {
	lineSummary := resultLineSummary{
		Component:   artifact.General.ComponentID,
		FixVersions: make(map[string]bool),
	}
	for _, issue := range artifact.Issues {
		switch issue.Severity {
		case "Medium":
			lineSummary.MediumCount++
		case "Low":
			lineSummary.LowCount++
		case "High":
			lineSummary.HighCount++
		}
		for _, component := range issue.Components {
			for _, version := range component.FixedVersions {
				lineSummary.FixVersions[version] = true
			}
		}
	}
	for _, license := range artifact.Licenses {
		lineSummary.Licenses = append(lineSummary.Licenses, license.Name)
	}
	return lineSummary
}

func versions(versions map[string]bool) string {
	var s []string
	for key := range versions {
		s = append(s, key)
	}
	return strings.Join(s, ", ")
}
