// 	jira-branch-helper - Build a string that can be used for a branch name from the details in a Jira ticket
//
// 	Copyright (C) 2017 Billie Alice Thompson
//
// 	This program is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU General Public License as published by
// 	the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
// 	This program is distributed in the hope that it will be useful,
// 	but WITHOUT ANY WARRANTY; without even the implied warranty of
// 	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 	GNU General Public License for more details.
//
// 	You should have received a copy of the GNU General Public License
// 	along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/PurpleBooth/jira-branch-helper/jira/branchhelper"
	"github.com/andygrunwald/go-jira"
	"github.com/urfave/cli"
	"strings"
)

const (
	// ErrorExitCodeIncorrectNumberOfArguments is the exit code when there isn't an issue provided
	ErrorExitCodeIncorrectNumberOfArguments = 1 << iota
	// ErrorExitCodeNoEndpointUrl is the exit code when there is no Jira endpoint
	ErrorExitCodeNoEndpointUrl
	// ErrorExitCodeJiraInitFailure happens when we fail to initialise our jira client
	ErrorExitCodeJiraInitFailure
	// ErrorExitCodeBranchNameBuildFailure happens when we fail to build the branch name
	ErrorExitCodeBranchNameBuildFailure
	// ErrorExitCodeCouldNotParseIssue happens when we fail to extract the issue number from what the user provided
	ErrorExitCodeCouldNotParseIssue
)

const ArgumentJiraEndpoint = "jira-endpoint"
const DefaultTemplate = "{{.Key | ToLower }}-{{.Fields.Summary | Trim | KebabCase }}"

var AppVersion string

func main() {
	app := cli.NewApp()
	app.Version = AppVersion
	app.Authors = []cli.Author{
		{
			Name:  "Billie Alice Thompson",
			Email: "billie@purplebooth.co.uk",
		},
	}
	app.Name = "jira-branch-helper"
	app.Usage = `
	Build a string that can be used for a branch name from the details in a Jira ticket

	Example usage

	$ jira-branch-helper https://example.com/jira/browse/TST-123
	tst-123-ticket-title-goes-here

	$ jira-branch-helper TST-123
	tst-123-ticket-title-goes-here

	Environment variables may be used in place of flags, parameters, see parameters with [$ENV_NAME_HERE] at the end.

	The following functions are available for templating

	* "Trim" - Remove whitespace at the start and end of the function (no parameters)
	* "ToLower" - Lower case your string (no parameters)
	* "ToUpper" - Upper case your string (no parameters)
	* "Replace" - Replace characters e.g. {{.Fields.Summary | Replace "A" "B" }}, would replace A with B
	* "KebabCase" - e.g. developments-phase-1-implement-feature--bume (no parameters)
	* "LowerSnakeCase" - e.g. developments_phase_1_implement_feature__bume (no parameters)
	* "LowerCamelCase" - e.g. developmentsPhase1ImplementFeatureBume (no parameters)
	* "ScreamingKebabCase" - e.g. DEVELOPMENTS-PHASE-1-IMPLEMENT-FEATURE--BUME (no parameters)
	* "ScreamingSnakeCase" - e.g. DEVELOPMENTS_PHASE_1_IMPLEMENT_FEATURE__BUME (no parameters)
	* "UpperCamelCase" - e.g. DevelopmentsPhase1ImplementFeatureBume (no parameters)
	* "UpperKebabCase" - e.g. Developments_Phase_1_Implement_Feature__Bume (no parameters)

	The template format is as described here https://golang.org/pkg/text/template/
	`
	app.Copyright = `
	jira-branch-helper  Copyright (C) 2017  Billie Alice Thompson
	This program comes with ABSOLUTELY NO WARRANTY;	This is free software,
	and you are welcome to redistribute it under certain conditions; see
	LICENSE.md for additional details.
	`

	app.ArgsUsage = "[ISSUE-NUMBER OR ISSUE-URL]"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_USERNAME",
			Name:   "jira-username",
			Usage:  "The username to authenticate as on Jira",
		},
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_PASSWORD",
			Name:   "jira-password",
			Usage:  "The password to authenticate as on Jira",
		},
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_ENDPOINT",
			Name:   ArgumentJiraEndpoint,
			Usage:  "Jira's URL",
		},
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_TEMPLATE",
			Name:   "template",
			Usage:  "The template to use to generate the branch name",
			Value:  DefaultTemplate,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError(
				"incorrect number of arguments, see `jira-branch-helper help` for full usage information",
				ErrorExitCodeIncorrectNumberOfArguments,
			)
		}
		var endpointUrl string
		var issueStrategy branchhelper.IssueStrategy
		unparsedIssue := c.Args().Get(0)
		issueUrl, err := url.Parse(unparsedIssue)

		if issueUrl != nil && issueUrl.Host != "" {
			issueStrategy = branchhelper.IssueUrlStrategy{}
		} else {
			issueStrategy = branchhelper.IssueLiteralStrategy{}
		}

		if c.String(ArgumentJiraEndpoint) == "" {
			urlHelpers := []branchhelper.EndpointStrategy{
				branchhelper.EndpointCombinedStrategy{},
				branchhelper.EndpointSoloStrategy{},
			}

			if err == nil {
				for i := range urlHelpers {
					possibleEndpointUrl := urlHelpers[i].GetEndpoint(issueUrl)

					if possibleEndpointUrl != nil {
						endpointUrl = possibleEndpointUrl.String()

						break
					}
				}
			}

			if endpointUrl == "" {
				return cli.NewExitError(
					"you must provide a Jira URL via Flag or environment variable or a full issue url",
					ErrorExitCodeNoEndpointUrl,
				)
			}
		} else {
			endpointUrl = c.String(ArgumentJiraEndpoint)

			if endpointUrl[len(endpointUrl)-1:] != "/" {
				endpointUrl = strings.Join([]string{endpointUrl, "/"}, "")
			}
		}

		jiraClient, err := jira.NewClient(nil, endpointUrl)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("failed to initialise jira client: %s", err.Error()),
				ErrorExitCodeJiraInitFailure,
			)
		}

		if c.String("jira-username") != "" {
			jiraClient.Authentication.SetBasicAuth(
				c.String("jira-username"),
				c.String("jira-password"),
			)
		}

		issueFormatter := branchhelper.NewBranchHelper(jiraClient)
		issueId, err := issueStrategy.GetIssue(unparsedIssue)

		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("failed parse unparsedIssue: %s", err.Error()),
				ErrorExitCodeCouldNotParseIssue,
			)
		}

		template := c.String("template")

		if template == "" {
			template = DefaultTemplate
		}

		branchName, err := issueFormatter.FormatIssue(
			issueId,
			template,
		)

		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("failed build branch name: %s", err.Error()),
				ErrorExitCodeBranchNameBuildFailure,
			)
		}

		os.Stdout.WriteString(branchName)
		os.Stdout.WriteString("\n")

		return nil
	}
	app.EnableBashCompletion = true
	app.Run(os.Args)
}
