package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBuildPrimarySection(t *testing.T) {
	var defaultValuesConfig = config{
		SectionTitle:                 "Some author",
		SectionSubtitle:              "A commit message",
		SectionText:                  "The commits message body",
		EnablePrimarySectionMarkdown: false,
	}
	var tests = []struct {
		input    config
		expected Section
	}{
		{
			defaultValuesConfig,
			Section{
				ActivityTitle:    defaultValuesConfig.SectionTitle,
				ActivitySubtitle: defaultValuesConfig.SectionSubtitle,
				Text:             defaultValuesConfig.SectionText,
				Markdown:         defaultValuesConfig.EnablePrimarySectionMarkdown,
			},
		},
	}

	for _, test := range tests {
		if output := buildPrimarySection(test.input); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Test failed: config input was %v, expected %v", test.input, test.expected)
		}
	}
}

func TestBuildFactsSection(t *testing.T) {
	const unixTimeString = "1610768692"
	const parsedUnixTime = "Sat, 16 Jan 2021 14:44:52 AEDT"

	var mockConfig = config{
		SuccessThemeColor:        "FFFFFF",
		FailedThemeColor:         "000000",
		EnableBuildFactsMarkdown: false,
		BuildNumber:              "1",
		BuildTime:                unixTimeString,
		GitBranch:                "master",
		Workflow:                 "master_branch",
	}

	var tests = []struct {
		input          config
		isBuildSuccess bool
		expected       Section
	}{
		// Successful build
		{
			mockConfig,
			true,
			Section{
				Markdown: mockConfig.EnableBuildFactsMarkdown,
				Facts: []Fact{
					{
						Name:  "Build Status",
						Value: fmt.Sprintf(`<span style="color:#%s">Success</span>`, mockConfig.SuccessThemeColor),
					},
					{
						Name:  "Build Number",
						Value: mockConfig.BuildNumber,
					},
					{
						Name:  "Git Branch",
						Value: mockConfig.GitBranch,
					},
					{
						Name:  "Build Triggered",
						Value: parsedUnixTime,
					},
					{
						Name:  "Workflow",
						Value: mockConfig.Workflow,
					},
				},
			},
		},
		// Failed build
		{
			mockConfig,
			false,
			Section{
				Markdown: mockConfig.EnableBuildFactsMarkdown,
				Facts: []Fact{
					{
						Name:  "Build Status",
						Value: fmt.Sprintf(`<span style="color:#%s">Fail</span>`, mockConfig.FailedThemeColor),
					},
					{
						Name:  "Build Number",
						Value: mockConfig.BuildNumber,
					},
					{
						Name:  "Git Branch",
						Value: mockConfig.GitBranch,
					},
					{
						Name:  "Build Triggered",
						Value: parsedUnixTime,
					},
					{
						Name:  "Workflow",
						Value: mockConfig.Workflow,
					},
				},
			},
		},
	}
	for _, test := range tests {
		if output := buildFactsSection(test.input, test.isBuildSuccess); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Test failed: config input was %v, expected %v", test.input, test.expected)
		}
	}
}
