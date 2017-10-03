// jira-branch-helper - Build a string that can be used for a branch name from
// the details in a Jira ticket
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
	"net/url"
	"os"
	"strings"

	"github.com/PurpleBooth/jira-branch-helper/jira/branchhelper"
	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	errorExitCodeIncorrectNumberOfArguments = 1 << iota
	errorExitCodeNoEndpointURL
	errorExitCodeJiraInitFailure
	errorExitCodeBranchNameBuildFailure
	errorExitCodeCouldNotParseIssue
	errorExitCodeBranchNameWriteError
)

const (
	// argumentJiraBasicUsername is the option to set the username for basic
	// auth for the Jira API
	argumentJiraBasicUsername = "jira-basic-auth-username"
	// argumentJiraBasicPassword is the option to set the username for basic
	// auth for the Jira API
	argumentJiraBasicPassword = "jira-basic-auth-password"
	// argumentJiraCookieUsername is the option to set the username for the Jira
	// API (via simulating a real login)
	argumentJiraCookieUsername = "jira-username"
	// argumentJiraCookiePassword is the option to set the password for the Jira
	// API (via simulating a real login)
	argumentJiraCookiePassword = "jira-password"
	// argumentJiraEndpoint is the option to set jira's URL
	argumentJiraEndpoint = "jira-endpoint"
	// argumentTemplate is the option to set the template to generate the branch
	argumentTemplate = "template"
)

// defaultTemplate is The default template to use for the branch
const defaultTemplate = "{{.Key | ToLower }}-" +
	"{{.Fields.Summary | Trim | KebabCase }}"

// AppVersion is the version of the app. Set when building using ldflags
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
	Build a string that can be used for a branch name from the details in a Jira
	ticket

	Example usage

	$ jira-branch-helper https://example.com/jira/browse/TST-123
	tst-123-ticket-title-goes-here

	$ jira-branch-helper TST-123
	tst-123-ticket-title-goes-here

	Environment variables may be used in place of flags, parameters, see
	parameters with [$ENV_NAME_HERE] at the end.

	The following functions are available for templating

	* "Trim"               - Remove whitespace from start and end
	* "ToLower"            - Lower case your string
	* "ToUpper"            - Upper case your string
	* "Replace"            - Replace characters params: for-search, replace-with
	* "KebabCase"          - Switch the casing-to-kebab
	* "LowerSnakeCase"     - Switch the casing_to_snake
	* "LowerCamelCase"     - Switch the casingToCamel
	* "ScreamingKebabCase" - Switch the CASING-TO-KEBAB
	* "ScreamingSnakeCase" - Switch the CASING_TO_SNAKE
	* "UpperCamelCase"     - Switch the CasingToCamel
	* "UpperKebabCase"     - Switch the Casing-To-Kebab

	The template format is as described here
	https://golang.org/pkg/text/template/

	Templates look like this

	{{.Key | ToLower }}-{{.Fields.Summary | Replace "A" "B" | KebabCase }}
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
			EnvVar: "JIRA_BRANCH_HELPER_USERNAME_BASIC_AUTH",
			Name:   argumentJiraBasicUsername,
			Usage:  "Set a basic auth username on HTTP requests to Jira",
		},
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_PASSWORD_BASIC_AUTH",
			Name:   argumentJiraBasicPassword,
			Usage:  "Set a basic auth password on HTTP requests to Jira",
		},
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_USERNAME",
			Name:   argumentJiraCookieUsername,
			Usage:  "The username to authenticate as on Jira",
		},
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_PASSWORD",
			Name:   argumentJiraCookiePassword,
			Usage:  "The password to authenticate as on Jira",
		},
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_ENDPOINT",
			Name:   argumentJiraEndpoint,
			Usage:  "Jira's URL",
		},
		cli.StringFlag{
			EnvVar: "JIRA_BRANCH_HELPER_TEMPLATE",
			Name:   argumentTemplate,
			Usage:  "The template to use to generate the branch name",
			Value:  defaultTemplate,
		},
	}
	app.Action = action
	app.EnableBashCompletion = true

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func action(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError(
			"incorrect number of arguments, see "+
				"`jira-branch-helper help` for full usage information",
			errorExitCodeIncorrectNumberOfArguments,
		)
	}

	var endpointURL string
	rawIssueID := c.Args().Get(0)
	issueURL, err := url.Parse(rawIssueID)
	issueStrategy := branchhelper.MakeIssueStrategy(issueURL)

	if c.String(argumentJiraEndpoint) == "" {
		endpointURL = branchhelper.GuessEndpointURL(issueURL)
		if endpointURL == "" {
			return cli.NewExitError(
				"you must provide a Jira URL via Flag or "+
					"environment variable or a full issue url",
				errorExitCodeNoEndpointURL,
			)
		}
	} else {
		endpointURL = c.String(argumentJiraEndpoint)
		endpointURL = normaliseEndpointURL(endpointURL)
	}

	jiraClient, err := jira.NewClient(nil, endpointURL)

	if err != nil {
		return cli.NewExitError(
			errors.Wrap(
				err,
				"initialising jira client failed",
			).Error(),
			errorExitCodeJiraInitFailure,
		)
	}

	if err := addSessionCookie(c, jiraClient); err != nil {
		return err
	}

	addBasicAuth(c, jiraClient)

	issueFormatter := branchhelper.NewJira(jiraClient)
	issueID, err := issueStrategy.GetIssue(rawIssueID)

	if err != nil {
		return cli.NewExitError(
			errors.Wrap(err, "failed to parse the issue id").Error(),
			errorExitCodeCouldNotParseIssue,
		)
	}

	template := c.String(argumentTemplate)
	branchName, err := formatIssue(issueFormatter, template, issueID)

	if err != nil {
		return cli.NewExitError(
			errors.Wrap(err, "failed to build branch name").Error(),
			errorExitCodeBranchNameBuildFailure,
		)
	} else {
		if _, err := os.Stdout.WriteString(branchName + "\n"); err != nil {
			wrappedErr := errors.Wrap(
				err,
				"failed to flush branch name to buffer",
			)

			return cli.NewExitError(
				wrappedErr.Error(),
				errorExitCodeBranchNameWriteError,
			)
		}
	}

	return nil
}
func addSessionCookie(c *cli.Context, jiraClient *jira.Client) *cli.ExitError {
	if c.String(argumentJiraCookieUsername) != "" {
		if _, err := jiraClient.Authentication.AcquireSessionCookie(
			c.String(argumentJiraCookieUsername),
			c.String(argumentJiraCookiePassword),
		); err != nil {
			wrappedErr := errors.Wrap(
				err,
				"failed to authenticate with jira",
			)

			return cli.NewExitError(
				wrappedErr.Error(),
				errorExitCodeJiraInitFailure,
			)
		}
	}

	return nil
}
func addBasicAuth(c *cli.Context, jiraClient *jira.Client) {
	if c.String(argumentJiraBasicUsername) != "" {
		jiraClient.Authentication.SetBasicAuth(
			c.String(argumentJiraBasicUsername),
			c.String(argumentJiraBasicPassword),
		)
	}
}

func normaliseEndpointURL(endpointURL string) string {
	if endpointURL[len(endpointURL)-1:] != "/" {
		endpointURL = strings.Join([]string{endpointURL, "/"}, "")
	}
	return endpointURL
}

func formatIssue(
	issueFormatter *branchhelper.Jira,
	template string,
	issueID string,
) (string, error) {

	if template == "" {
		template = defaultTemplate
	}
	branchName, err := issueFormatter.FormatIssue(
		issueID,
		template,
	)
	return branchName, err
}
