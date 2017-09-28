package branchhelper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/PurpleBooth/jira-branch-helper/jira/branchhelper"
	"net/url"
)

var _ = Describe("EndpointSoloStrategy", func() {
	Context("Failure", func() {
		It("Returns empty template", func() {
			actual := EndpointSoloStrategy{}
			input, _ := url.Parse("https://example.com/not-a-issue")

			Expect(actual.GetEndpoint(input)).To(BeNil())
		})
		It("Fails on non-solo url", func() {
			actual := EndpointSoloStrategy{}
			input, _ := url.Parse("https://example.com/jira/browse/TST-123")

			Expect(actual.GetEndpoint(input)).To(BeNil())
		})
	})
	Context("Success", func() {
		It("Returns empty template", func() {
			actual := EndpointSoloStrategy{}
			input, _ := url.Parse("https://example.com/browse/TST-123")
			expected, _ := url.Parse("https://example.com/")

			Expect(actual.GetEndpoint(input)).To(Equal(expected))
		})
	})
})

var _ = Describe("EndpointCombinedStrategy", func() {
	Context("Failure", func() {
		It("Fails on non-issue url", func() {
			actual := EndpointCombinedStrategy{}
			input, _ := url.Parse("https://example.com/not-a-issue")

			Expect(actual.GetEndpoint(input)).To(BeNil())
		})
		It("Fails on non-combined url", func() {
			actual := EndpointCombinedStrategy{}
			input, _ := url.Parse("https://example.com/browse/TST-123")

			Expect(actual.GetEndpoint(input)).To(BeNil())
		})
	})
	Context("Success", func() {
		It("Returns empty template", func() {
			actual := EndpointCombinedStrategy{}
			input, _ := url.Parse("https://example.com/jira/browse/TST-123")
			expected, _ := url.Parse("https://example.com/jira")

			Expect(actual.GetEndpoint(input)).To(Equal(expected))
		})
	})
})
