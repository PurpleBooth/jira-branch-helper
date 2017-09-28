package branchhelper

import (
	"bufio"
	"bytes"
	"strings"
	"text/template"

	"regexp"

	"github.com/andygrunwald/go-jira"
	"github.com/danverbraganza/varcaser/varcaser"
)

type BranchHelper struct {
	Client GetIssueClient
}

type GetIssueClient interface {
	Get(issueID string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error)
}

func toSnakeCase(s string) string {
	unneededCharactersReg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	lowWithSpace := strings.ToLower(unneededCharactersReg.ReplaceAllString(s, ""))
	spacesReg, _ := regexp.Compile("[ ]")
	return spacesReg.ReplaceAllString(lowWithSpace, "_")
}

func trim(s string) string {
	startRegex, _ := regexp.Compile("^\\s+")
	endRegex, _ := regexp.Compile("\\s+$")
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

func toVarCaserFunction(convention varcaser.CaseConvention) func(string) string {
	return toSnakeCaseFunction(varcaser.Caser{From: varcaser.LowerSnakeCase, To: convention}.String)
}

func (helper *BranchHelper) FormatIssue(issueId string, rawTmpl string) (string, error) {
	templ, err := template.New("branch-name").Funcs(template.FuncMap{
		"Trim":               trim,
		"ToLower":            strings.ToLower,
		"ToUpper":            strings.ToUpper,
		"Replace":            replace,
		"LowerSnakeCase":     toSnakeCase,
		"KebabCase":          toVarCaserFunction(varcaser.KebabCase),
		"LowerCamelCase":     toVarCaserFunction(varcaser.LowerCamelCase),
		"ScreamingKebabCase": toVarCaserFunction(varcaser.ScreamingKebabCase),
		"ScreamingSnakeCase": toVarCaserFunction(varcaser.ScreamingSnakeCase),
		"UpperCamelCase":     toVarCaserFunction(varcaser.UpperCamelCase),
		"UpperKebabCase":     toVarCaserFunction(varcaser.UpperKebabCase),
	}).Parse(rawTmpl)

	if err != nil {
		return "", err
	}

	issue, _, err := helper.Client.Get(issueId, nil)
	output := &bytes.Buffer{}
	writer := bufio.NewWriter(output)

	if err != nil {
		return "", err
	}

	err = templ.Execute(writer, issue)
	writer.Flush()

	if err != nil {
		return "", err
	}

	return output.String(), nil

}

func NewBranchHelper(client *jira.Client) *BranchHelper {
	return &BranchHelper{Client: client.Issue}
}
