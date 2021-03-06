package main

import (
	"fmt"
	"reflect"
	"testing"
)

const unixTimeString = "1610768692"
const parsedUnixTime = "Sat, 16 Jan 2021 03:44:52 UTC"

var mockConfig = config{
	BuildNumber:                  "1",
	AppTitle:                     "Some app title",
	AppURL:                       "https://www.github.com/username/repo",
	BuildURL:                     "https://www.bitrise.io/some/app/build",
	BuildTime:                    unixTimeString,
	CommitMessageBody:            "Some commit message body",
	GitBranch:                    "master",
	Workflow:                     "master_branch",
	WebhookURL:                   "https://microsoft.com/some/webhook",
	CardTitle:                    "The heading for the card",
	SuccessThemeColor:            "FFFFFF",
	FailedThemeColor:             "000000",
	SectionTitle:                 "Git author name",
	SectionSubtitle:              "Commit message",
	SectionText:                  "Commit message body",
	SectionHeaderImage:           "",
	EnablePrimarySectionMarkdown: "no",
	EnableBuildFactsMarkdown:     "no",
	EnableDefaultActions:         "yes",
	EnableDebug:                  "no",
	RepoURL:                      "https://www.github.com/username/repo",
	Actions: `[
		{
			"text": "Some text",
			"targets": [
				{
					"uri": "www.google.com", 
					"os": "android"
				},
				{
					"uri": "www.google.com", 
					"os": "iOS"
				},
				{
					"uri": "www.google.com", 
					"os": "windows"
				}
			]
		}
	]`,
}

func TestGetValueForBuildStatus(t *testing.T) {
	const success = "success"
	const fail = "fail"

	successValue := getValueForBuildStatus(success, fail, true)
	failValue := getValueForBuildStatus(success, fail, false)

	if successValue != success {
		t.Errorf("Test failed: expected %v but input was %v", success, successValue)
	}
	if failValue != fail {
		t.Errorf("Test failed: expected %v but input was %v", fail, failValue)
	}
}

func TestOptionalUserValue(t *testing.T) {
	const defaultValue = "default value"
	const userValue = "user value"

	fallbackToDefault := optionalUserValue(defaultValue, "")
	userCustomValue := optionalUserValue(defaultValue, userValue)
	if fallbackToDefault != defaultValue {
		t.Errorf("Test failed: expected %v but input was %v", defaultValue, fallbackToDefault)
	}
	if userCustomValue != userValue {
		t.Errorf("Test failed: expected %v but input was %v", userValue, userCustomValue)
	}
}

func TestBuildPrimarySection(t *testing.T) {
	var defaultValuesConfig = config{
		SectionTitle:                 "Some author",
		SectionSubtitle:              "A commit message",
		SectionText:                  "The commits message body",
		EnablePrimarySectionMarkdown: "no",
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
				Markdown:         valueOptionToBool(defaultValuesConfig.EnablePrimarySectionMarkdown),
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
				Markdown: valueOptionToBool(mockConfig.EnableBuildFactsMarkdown),
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
				Markdown: valueOptionToBool(mockConfig.EnableBuildFactsMarkdown),
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

func TestNewMessage(t *testing.T) {

	var buildSuccessFacts = Section{
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
		Markdown:  valueOptionToBool(mockConfig.EnableBuildFactsMarkdown),
		HeroImage: HeroImage{},
	}

	var buildFailedFacts = Section{
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
		Markdown:  valueOptionToBool(mockConfig.EnableBuildFactsMarkdown),
		HeroImage: HeroImage{},
	}

	var primarySection = Section{
		ActivityTitle:    mockConfig.SectionTitle,
		ActivitySubtitle: mockConfig.SectionSubtitle,
		ActivityImage:    mockConfig.SectionHeaderImage,
		Markdown:         valueOptionToBool(mockConfig.EnablePrimarySectionMarkdown),
		Text:             mockConfig.SectionText,
		HeroImage:        HeroImage{},
	}

	var buildSuccessMessage = Message{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: mockConfig.SuccessThemeColor,
		Title:      mockConfig.CardTitle,
		Summary:    fmt.Sprintf("%v #%v succeeded", mockConfig.AppTitle, mockConfig.BuildNumber),
		// Be mindful of list order
		Sections: []Section{primarySection, buildSuccessFacts},
		Actions: []OpenURIAction{
			{
				Type: "OpenUri",
				Name: "Go To Repo",
				Targets: []Target{
					{
						OS:  "default",
						URI: mockConfig.RepoURL,
					},
				},
			},
			{
				Type: "OpenUri",
				Name: "Go To Build",
				Targets: []Target{
					{
						OS:  "default",
						URI: mockConfig.BuildURL,
					},
				},
			},
			{
				Type: "OpenUri",
				Name: "Some text",
				Targets: []Target{
					{
						URI: "www.google.com",
						OS:  "android",
					},
					{
						URI: "www.google.com",
						OS:  "iOS",
					},
					{
						URI: "www.google.com",
						OS:  "windows",
					},
				},
			},
		},
	}

	var buildFailedMessage = Message{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: mockConfig.FailedThemeColor,
		Title:      mockConfig.CardTitle,
		Summary:    fmt.Sprintf("%v #%v failed", mockConfig.AppTitle, mockConfig.BuildNumber),
		// Be mindful of list order
		Sections: []Section{primarySection, buildFailedFacts},
		Actions: []OpenURIAction{
			{
				Type: "OpenUri",
				Name: "Go To Repo",
				Targets: []Target{
					{
						OS:  "default",
						URI: mockConfig.RepoURL,
					},
				},
			},
			{
				Type: "OpenUri",
				Name: "Go To Build",
				Targets: []Target{
					{
						OS:  "default",
						URI: mockConfig.BuildURL,
					},
				},
			},
			{
				Type: "OpenUri",
				Name: "Some text",
				Targets: []Target{
					{
						URI: "www.google.com",
						OS:  "android",
					},
					{
						URI: "www.google.com",
						OS:  "iOS",
					},
					{
						URI: "www.google.com",
						OS:  "windows",
					},
				},
			},
		},
	}

	var tests = []struct {
		input          config
		isBuildSuccess bool
		expected       Message
	}{
		{
			mockConfig,
			true,
			buildSuccessMessage,
		},
		{
			mockConfig,
			false,
			buildFailedMessage,
		},
	}
	for _, test := range tests {
		if output := newMessage(test.input, test.isBuildSuccess); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Test failed: config input was %v, expected %v", test.input, test.expected)
		}
	}

}

func TestBuildURIAction(t *testing.T) {
	var tests = []struct {
		input    Action
		expected OpenURIAction
	}{
		{
			input: Action{
				Text: "Action 1",
				Targets: []ActionTarget{
					{
						URI: "https://www.google.com",
					},
				},
			},
			expected: OpenURIAction{
				Type: "OpenUri",
				Name: "Action 1",
				Targets: []Target{
					{
						OS:  "default",
						URI: "https://www.google.com",
					},
				},
			},
		},
		{
			input: Action{
				Text: "Action 2",
				Targets: []ActionTarget{
					{
						URI: "https://www.google.com",
						OS:  "iOS",
					},
					{
						URI: "https://www.google.com",
						OS:  "android",
					},
					{
						URI: "https://www.google.com",
						OS:  "windows",
					},
					{
						URI: "https://www.google.com",
						OS:  "default",
					},
				},
			},
			expected: OpenURIAction{
				Type: "OpenUri",
				Name: "Action 2",
				Targets: []Target{
					{
						OS:  "iOS",
						URI: "https://www.google.com",
					},
					{
						OS:  "android",
						URI: "https://www.google.com",
					},
					{
						OS:  "windows",
						URI: "https://www.google.com",
					},
					{
						OS:  "default",
						URI: "https://www.google.com",
					},
				},
			},
		},
	}

	for _, test := range tests {
		if output := buildURIAction(test.input); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("Test failed: config input was %v, expected %v", test.input, test.expected)
		}
	}
}
