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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("IssueLiteralStrategy", func() {
	Context("Success", func() {
		It("Does nothing to urls ", func() {
			actual := IssueLiteralStrategy{}

			Expect(
				actual.GetIssue(
					"https://example.com/browse/TST-123",
				),
			).To(
				Equal(
					"https://example.com/browse/TST-123",
				),
			)
		})
		It("Excepts numbers ", func() {
			actual := IssueLiteralStrategy{}

			Expect(
				actual.GetIssue("TST-123"),
			).To(
				Equal("TST-123"),
			)
		})
	})
})

var _ = Describe("IssueURLStrategy", func() {
	Context("Success", func() {
		It("Pulls the ID from URLS", func() {
			actual, err := (IssueURLStrategy{}).GetIssue(
				"https://example.com/browse/TST-123",
			)

			Expect(actual).To(Equal("TST-123"))
			Expect(err).To(BeNil())
		})
		It("Does not accept Ids", func() {
			tester := IssueURLStrategy{}
			actual, err := tester.GetIssue("Not an issue?")

			Expect(actual).To(Equal(""))
			Expect(err).ToNot(BeNil())
		})
	})
})

var _ = Describe("MakeIssueStrategy", func() {
	Context("literal strategy", func() {
		It("No url", func() {
			actual := MakeIssueStrategy(nil)

			Expect(actual).To(BeAssignableToTypeOf(IssueLiteralStrategy{}))
		})
		It("No host", func() {
			url, _ := url.Parse("https://")
			actual := MakeIssueStrategy(url)

			Expect(actual).To(BeAssignableToTypeOf(IssueLiteralStrategy{}))
		})
	})
	Context("url strategy", func() {
		It("No host", func() {
			url, _ := url.Parse("https://example.com/jira/TST-101")
			actual := MakeIssueStrategy(url)

			Expect(actual).To(BeAssignableToTypeOf(IssueURLStrategy{}))
		})
	})
})
