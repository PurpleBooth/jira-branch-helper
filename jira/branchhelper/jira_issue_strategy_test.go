package branchhelper_test

import (
	. "github.com/purplebooth/jira-branch-helper/jira/branchhelper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IssueLiteralStrategy", func() {
	Context("Success", func() {
		It("Does nothing to urls ", func() {
			actual := IssueLiteralStrategy{}

			Expect(actual.GetIssue("https://example.com/browse/TST-123")).To(Equal("https://example.com/browse/TST-123"))
		})
		It("Excepts numbers ", func() {
			actual := IssueLiteralStrategy{}

			Expect(actual.GetIssue("TST-123")).To(Equal("TST-123"))
		})
	})
})

var _ = Describe("IssueUrlStrategy", func() {
	Context("Success", func() {
		It("Pulls the ID from URLS", func() {
			actual, err := (IssueUrlStrategy{}).GetIssue("https://example.com/browse/TST-123")

			Expect(actual).To(Equal("TST-123"))
			Expect(err).To(BeNil())
		})
		It("Does not accept Ids", func() {
			tester := IssueUrlStrategy{}
			actual, err := tester.GetIssue("Not an issue?")

			Expect(actual).To(Equal(""))
			Expect(err).ToNot(BeNil())
		})
	})
})
