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

package branchhelper

import (
	"bufio"
	"bytes"
	"errors"
	"net/http/httputil"
	"regexp"
	"strings"
	"text/template"

	"github.com/andygrunwald/go-jira"
	"github.com/danverbraganza/varcaser/varcaser"
)

// Jira will generate branch names from Jira issues
type Jira struct {
	Client GetIssueClient
}

// GetIssueClient allows us to get issues from Jira
type GetIssueClient interface {
	Get(
		issueID string,
		options *jira.GetQueryOptions,
	) (
		*jira.Issue,
		*jira.Response,
		error,
	)
}

func toSnakeCase(s string) string {
	unneededCharactersReg, err := regexp.Compile("[^a-zA-Z0-9 ]+")

	if err != nil {
		return s
	}

	lowWithSpace := strings.ToLower(unneededCharactersReg.ReplaceAllString(s, ""))
	spacesReg, err := regexp.Compile("[ ]")

	if err != nil {
		return s
	}

	return spacesReg.ReplaceAllString(lowWithSpace, "_")
}

func trim(s string) string {
	startRegex, err := regexp.Compile("^\\s+")

	if err != nil {
		return s
	}

	endRegex, err := regexp.Compile("\\s+$")

	if err != nil {
		return s
	}

	return endRegex.ReplaceAllString(startRegex.ReplaceAllString(s, ""), "")
}

func replace(replace string, with string, s string) string {
	return strings.Replace(s, replace, with, -1)
}

func toSnakeCaseFunction(stringFunc func(string) string) func(string) string {
	return func(s string) string {
		return stringFunc(toSnakeCase(s))
	}
}

func normaliseArgument(
	convention varcaser.CaseConvention,
) func(string) string {
	return toSnakeCaseFunction(
		varcaser.Caser{
			From: varcaser.LowerSnakeCase,
			To:   convention,
		}.String,
	)
}

// FormatIssue generate a branch name from a template and a issue ID
func (helper *Jira) FormatIssue(
	issueID string,
	rawTmpl string,
) (string, error) {
	templ, err := template.New("branch-name").Funcs(template.FuncMap{
		"Trim":               trim,
		"ToLower":            strings.ToLower,
		"ToUpper":            strings.ToUpper,
		"Replace":            replace,
		"LowerSnakeCase":     toSnakeCase,
		"KebabCase":          normaliseArgument(varcaser.KebabCase),
		"LowerCamelCase":     normaliseArgument(varcaser.LowerCamelCase),
		"ScreamingKebabCase": normaliseArgument(varcaser.ScreamingKebabCase),
		"ScreamingSnakeCase": normaliseArgument(varcaser.ScreamingSnakeCase),
		"UpperCamelCase":     normaliseArgument(varcaser.UpperCamelCase),
		"UpperKebabCase":     normaliseArgument(varcaser.UpperKebabCase),
	}).Parse(rawTmpl)

	if err != nil {
		return "", err
	}

	issue, resp, err := helper.Client.Get(issueID, nil)
	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)

	if err != nil {
		return "", makeJiraRequestError(err, resp)
	}

	if err := templ.Execute(writer, issue); err != nil {
		return "", err
	}

	if err := writer.Flush(); err != nil {
		return "", err
	}

	return buffer.String(), nil

}

func makeJiraRequestError(requestErr error, resp *jira.Response) error {
	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)

	if resp != nil {
		dumpedResponse, dumpErr := jiraResponseToString(resp)

		if dumpErr != nil {
			return dumpErr
		}

		if _, err := buffer.WriteString(dumpedResponse); err != nil {
			return err
		}
	}

	if _, err := writer.WriteString(requestErr.Error()); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return errors.New(buffer.String())
}

func jiraResponseToString(resp *jira.Response) (string, error) {
	buf := &bytes.Buffer{}
	writer := bufio.NewWriter(buf)

	reqDump, err := httputil.DumpRequest(resp.Request, true)

	if err != nil {
		return "", err
	}

	if _, err = writer.Write(reqDump); err != nil {
		return "", err
	}

	if _, err = writer.WriteString("\n\n"); err != nil {
		return "", err
	}

	respDump, err := httputil.DumpResponse(resp.Response, true)

	if err != nil {
		return "", err
	}

	if _, err = writer.Write(respDump); err != nil {
		return "", err
	}

	if _, err = writer.WriteString("\n\n"); err != nil {
		return "", err
	}

	return buf.String(), err
}

// NewJira from a Jira client build a client helper
func NewJira(
	client *jira.Client,
) *Jira {
	return &Jira{Client: client.Issue}
}
