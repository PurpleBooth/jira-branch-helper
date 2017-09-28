package branchhelper_test

import (
	"github.com/andygrunwald/go-jira"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/purplebooth/jira-branch-helper/jira/branchhelper"
)

var _ = Describe("BranchHelper", func() {
	Context("Templating", func() {
		It("Error in template causes error", func() {
			subject := BranchHelper{}
			actual, err := subject.FormatIssue("TST-123", "{{ .I am a bo} sdfsdfeef {{ .Broken }}")

			Expect(actual).To(Equal(""))
			Expect(err).ToNot(BeNil())
		})
		It("Template with no functions in has access to issue still", func() {
			subject := BranchHelper{Client: testGetIssue{
				issue:    &jira.Issue{Key: "TST-123"},
				response: nil,
				err:      nil,
			}}

			actual, err := subject.FormatIssue("TST-123", "{{.Key }}")

			Expect(actual).To(Equal("TST-123"))
			Expect(err).To(BeNil())
		})
		It("Has trim function", func() {
			subject := BranchHelper{Client: testGetIssue{
				issue: &jira.Issue{Fields: &jira.IssueFields{
					Summary: "    Developments Phase 1: Implement Feature γ Bäume    ",
				}},
				response: nil,
				err:      nil,
			}}

			actual, err := subject.FormatIssue("TST-123", "{{ .Fields.Summary | Trim }}")

			Expect(actual).To(Equal("Developments Phase 1: Implement Feature γ Bäume"))
			Expect(err).To(BeNil())
		})
		It("Can lower case", func() {
			subject := BranchHelper{Client: testGetIssue{
				issue: &jira.Issue{Fields: &jira.IssueFields{
					Summary: "Developments Phase 1: Implement Feature γ Bäume",
				}},
				response: nil,
				err:      nil,
			}}

			actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | ToLower }}")

			Expect(actual).To(Equal("developments phase 1: implement feature γ bäume"))
			Expect(err).To(BeNil())
		})
		It("Can upper case", func() {
			subject := BranchHelper{Client: testGetIssue{
				issue: &jira.Issue{Fields: &jira.IssueFields{
					Summary: "Developments Phase 1: Implement Feature γ Bäume",
				}},
				response: nil,
				err:      nil,
			}}

			actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | ToUpper }}")

			Expect(actual).To(Equal("DEVELOPMENTS PHASE 1: IMPLEMENT FEATURE Γ BÄUME"))
			Expect(err).To(BeNil())
		})
		It("Can replace", func() {
			subject := BranchHelper{Client: testGetIssue{
				issue: &jira.Issue{Fields: &jira.IssueFields{
					Summary: "Developments Phase 1: Implement Feature γ Bäume",
				}},
				response: nil,
				err:      nil,
			}}

			actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | Replace \":\" \"!\" }}")

			Expect(actual).To(Equal("Developments Phase 1! Implement Feature γ Bäume"))
			Expect(err).To(BeNil())
		})
	})
	It("Lower snake case", func() {
		subject := BranchHelper{Client: testGetIssue{
			issue: &jira.Issue{Fields: &jira.IssueFields{
				Summary: "Developments Phase 1: Implement Feature γ Bäume",
			}},
			response: nil,
			err:      nil,
		}}

		actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | LowerSnakeCase }}")

		Expect(actual).To(Equal("developments_phase_1_implement_feature__bume"))
		Expect(err).To(BeNil())
	})
	It("KebabCase", func() {
		subject := BranchHelper{Client: testGetIssue{
			issue: &jira.Issue{Fields: &jira.IssueFields{
				Summary: "Developments Phase 1: Implement Feature γ Bäume",
			}},
			response: nil,
			err:      nil,
		}}

		actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | KebabCase }}")

		Expect(actual).To(Equal("developments-phase-1-implement-feature--bume"))
		Expect(err).To(BeNil())
	})
	It("LowerCamelCase", func() {
		subject := BranchHelper{Client: testGetIssue{
			issue: &jira.Issue{Fields: &jira.IssueFields{
				Summary: "Developments Phase 1: Implement Feature γ Bäume",
			}},
			response: nil,
			err:      nil,
		}}

		actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | LowerCamelCase }}")

		Expect(actual).To(Equal("developmentsPhase1ImplementFeatureBume"))
		Expect(err).To(BeNil())
	})
	It("ScreamingKebabCase", func() {
		subject := BranchHelper{Client: testGetIssue{
			issue: &jira.Issue{Fields: &jira.IssueFields{
				Summary: "Developments Phase 1: Implement Feature γ Bäume",
			}},
			response: nil,
			err:      nil,
		}}

		actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | ScreamingKebabCase }}")

		Expect(actual).To(Equal("DEVELOPMENTS-PHASE-1-IMPLEMENT-FEATURE--BUME"))
		Expect(err).To(BeNil())
	})
	It("ScreamingSnakeCase", func() {
		subject := BranchHelper{Client: testGetIssue{
			issue: &jira.Issue{Fields: &jira.IssueFields{
				Summary: "Developments Phase 1: Implement Feature γ Bäume",
			}},
			response: nil,
			err:      nil,
		}}

		actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | ScreamingSnakeCase }}")

		Expect(actual).To(Equal("DEVELOPMENTS_PHASE_1_IMPLEMENT_FEATURE__BUME"))
		Expect(err).To(BeNil())
	})
	It("UpperCamelCase", func() {
		subject := BranchHelper{Client: testGetIssue{
			issue: &jira.Issue{Fields: &jira.IssueFields{
				Summary: "Developments Phase 1: Implement Feature γ Bäume",
			}},
			response: nil,
			err:      nil,
		}}

		actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | UpperCamelCase }}")

		Expect(actual).To(Equal("DevelopmentsPhase1ImplementFeatureBume"))
		Expect(err).To(BeNil())
	})
	It("UpperKebabCase", func() {
		subject := BranchHelper{Client: testGetIssue{
			issue: &jira.Issue{Fields: &jira.IssueFields{
				Summary: "Developments Phase 1: Implement Feature γ Bäume",
			}},
			response: nil,
			err:      nil,
		}}

		actual, err := subject.FormatIssue("TST-123", "{{.Fields.Summary | UpperKebabCase }}")

		Expect(actual).To(Equal("Developments-Phase-1-Implement-Feature--Bume"))
		Expect(err).To(BeNil())
	})
})

type testGetIssue struct {
	issue    *jira.Issue
	response *jira.Response
	err      error
}

func (t testGetIssue) Get(issueID string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error) {
	return t.issue, t.response, t.err
}
