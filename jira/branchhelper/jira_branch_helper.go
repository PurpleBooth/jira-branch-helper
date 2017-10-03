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
	"net/http/httputil"
	"regexp"
	"strings"
	"text/template"

	"github.com/andygrunwald/go-jira"
	"github.com/danverbraganza/varcaser/varcaser"
	"github.com/pkg/errors"
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
	rawTempl string,
) (string, error) {
	templ, err := template.New(
		"branch-name",
	).Funcs(
		templateFunctions(),
	).Parse(rawTempl)

	if err != nil {
		return "", errors.Wrap(
			err,
			"failed to parse branch template",
		)
	}

	issue, resp, err := helper.Client.Get(issueID, nil)
	if err != nil {
		return "", newRequestError(err, resp)
	}

	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)

	if err := templ.Execute(writer, issue); err != nil {
		return "", errors.Wrap(
			err,
			"failed to execute branch template",
		)
	}

	if err := writer.Flush(); err != nil {
		return "", errors.Wrap(
			err,
			"failed flush template output to buffer",
		)
	}

	return buffer.String(), nil

}

func newRequestError(triggerErr error, resp *jira.Response) error {
	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)

	if resp != nil {
		respStr, err := dumpResponse(resp)

		if err != nil {
			return errors.Wrap(err, "failed dumping jira response")
		}

		if _, err := buffer.WriteString(respStr); err != nil {
			return errors.Wrap(
				err,
				"failed writing jira response to buffer",
			)
		}
	}

	if _, err := writer.WriteString(triggerErr.Error()); err != nil {
		return errors.Wrap(
			err,
			"failed writing jira error to buffer",
		)
	}

	if err := writer.Flush(); err != nil {
		return errors.Wrap(
			err,
			"failed flushing jira error to buffer",
		)
	}

	return errors.New(buffer.String())
}

func dumpResponse(resp *jira.Response) (string, error) {
	respParts := []string{}

	reqDump, err := httputil.DumpRequest(resp.Request, true)
	if err != nil {
		return "", errors.Wrap(
			err,
			"dumping jira http request failed",
		)
	}

	respDump, err := httputil.DumpResponse(resp.Response, true)
	if err != nil {
		return "", errors.Wrap(
			err,
			"dumping jira http response failed",
		)
	}

	respParts = append(respParts, string(reqDump))
	respParts = append(respParts, "\n\n")
	respParts = append(respParts, string(respDump))
	respParts = append(respParts, "\n\n")
	fullResp := strings.Join(respParts, "")

	return fullResp, nil
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
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
	}
}

// NewJira from a Jira client, build a client helper
func NewJira(
	client *jira.Client,
) *Jira {
	return &Jira{Client: client.Issue}
}
