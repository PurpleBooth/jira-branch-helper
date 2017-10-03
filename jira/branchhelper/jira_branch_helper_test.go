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

package branchhelper_test

import (
	. "github.com/PurpleBooth/jira-branch-helper/jira/branchhelper"
	"github.com/andygrunwald/go-jira"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jira", func() {
	Context("Templating", func() {
		It("Error in template causes error", func() {
			subject := Jira{}
			actual, err := subject.FormatIssue(
				"TST-123",
				"{{ .I am a bo} sdfsdfeef {{ .Broken }}",
			)

			Expect(actual).To(Equal(""))
			Expect(err).ToNot(BeNil())
		})
		It("Template with no functions in has access to issue still", func() {
			subject := Jira{Client: testGetIssue{
				issue:    &jira.Issue{Key: "TST-123"},
				response: nil,
				err:      nil,
			}}

			actual, err := subject.FormatIssue("TST-123", "{{.Key }}")

			Expect(actual).To(Equal("TST-123"))
			Expect(err).To(BeNil())
		})
		It("Has trim function", func() {
			actual, err := formatIssue(
				"    Developments Phase 1: Implement Feature γ Bäume    ",
				"{{.Fields.Summary | Trim }}",
			)

			Expect(actual).To(Equal("Developments Phase 1: Implement Feature γ Bäume"))
			Expect(err).To(BeNil())
		})
		It("Can lower case", func() {
			actual, err := formatIssue(
				"Developments Phase 1: Implement Feature γ Bäume",
				"{{.Fields.Summary | ToLower }}",
			)

			Expect(actual).To(Equal("developments phase 1: implement feature γ bäume"))
			Expect(err).To(BeNil())
		})
		It("Can upper case", func() {
			actual, err := formatIssue(
				"Developments Phase 1: Implement Feature γ Bäume",
				"{{.Fields.Summary | ToUpper }}",
			)

			Expect(actual).To(Equal("DEVELOPMENTS PHASE 1: IMPLEMENT FEATURE Γ BÄUME"))
			Expect(err).To(BeNil())
		})
		It("Can replace", func() {
			actual, err := formatIssue(
				"Developments Phase 1: Implement Feature γ Bäume",
				"{{.Fields.Summary | Replace \":\" \"!\" }}",
			)

			Expect(actual).To(Equal("Developments Phase 1! Implement Feature γ Bäume"))
			Expect(err).To(BeNil())
		})
	})
	It("Lower snake case", func() {
		actual, err := formatIssue(
			"Developments Phase 1: Implement Feature γ Bäume",
			"{{.Fields.Summary | LowerSnakeCase }}",
		)

		Expect(actual).To(Equal("developments_phase_1_implement_feature__bume"))
		Expect(err).To(BeNil())
	})
	It("KebabCase", func() {
		actual, err := formatIssue(
			"Developments Phase 1: Implement Feature γ Bäume",
			"{{.Fields.Summary | KebabCase }}",
		)

		Expect(actual).To(Equal("developments-phase-1-implement-feature--bume"))
		Expect(err).To(BeNil())
	})
	It("LowerCamelCase", func() {
		actual, err := formatIssue(
			"Developments Phase 1: Implement Feature γ Bäume",
			"{{.Fields.Summary | LowerCamelCase }}",
		)

		Expect(actual).To(Equal("developmentsPhase1ImplementFeatureBume"))
		Expect(err).To(BeNil())
	})
	It("ScreamingKebabCase", func() {
		actual, err := formatIssue(
			"Developments Phase 1: Implement Feature γ Bäume",
			"{{.Fields.Summary | ScreamingKebabCase }}",
		)

		Expect(actual).To(Equal("DEVELOPMENTS-PHASE-1-IMPLEMENT-FEATURE--BUME"))
		Expect(err).To(BeNil())
	})
	It("ScreamingSnakeCase", func() {
		actual, err := formatIssue(
			"Developments Phase 1: Implement Feature γ Bäume",
			"{{.Fields.Summary | ScreamingSnakeCase }}",
		)

		Expect(actual).To(Equal("DEVELOPMENTS_PHASE_1_IMPLEMENT_FEATURE__BUME"))
		Expect(err).To(BeNil())
	})
	It("UpperCamelCase", func() {
		actual, err := formatIssue(
			"Developments Phase 1: Implement Feature γ Bäume",
			"{{.Fields.Summary | UpperCamelCase }}",
		)

		Expect(actual).To(Equal("DevelopmentsPhase1ImplementFeatureBume"))
		Expect(err).To(BeNil())
	})
	It("UpperKebabCase", func() {
		actual, err := formatIssue(
			"Developments Phase 1: Implement Feature γ Bäume",
			"{{.Fields.Summary | UpperKebabCase }}",
		)

		Expect(actual).To(Equal("Developments-Phase-1-Implement-Feature--Bume"))
		Expect(err).To(BeNil())
	})
})

func formatIssue(summary string, templ string) (string, error) {
	subject := Jira{Client: testGetIssue{
		issue: &jira.Issue{Fields: &jira.IssueFields{
			Summary: summary,
		}},
		response: nil,
		err:      nil,
	}}

	return subject.FormatIssue(
		"TST-123",
		templ,
	)
}

type testGetIssue struct {
	issue    *jira.Issue
	response *jira.Response
	err      error
}

func (t testGetIssue) Get(
	issueID string,
	options *jira.GetQueryOptions,
) (*jira.Issue, *jira.Response, error) {
	return t.issue, t.response, t.err
}
