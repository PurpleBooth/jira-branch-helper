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
	"net/url"

	. "github.com/PurpleBooth/jira-branch-helper/jira/branchhelper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = Describe("Guess endpoint url", func() {
	Context("Failure", func() {
		It("Fails on non-issue url", func() {
			actual := GuessEndpointURL(nil)
			Expect(actual).To(Equal(""))
		})
	})
	Context("Success", func() {
		It("It uses combined strategy", func() {
			url, _ := url.Parse("https://example.com/jira/browse/TST-101")

			actual := GuessEndpointURL(url)
			Expect(actual).To(Equal("https://example.com/jira"))
		})
		It("It uses solo strategy", func() {
			url, _ := url.Parse("https://example.com/browse/TST-101")

			actual := GuessEndpointURL(url)
			Expect(actual).To(Equal("https://example.com/"))
		})
	})
})
