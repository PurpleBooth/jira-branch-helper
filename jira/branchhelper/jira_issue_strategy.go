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
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// IssueStrategy strategy to use to take a user supplied issue and convert it to
// a standard Jira issue
type IssueStrategy interface {
	GetIssue(rawIssue string) (string, error)
}

// IssueLiteralStrategy take what the user provided verboten
type IssueLiteralStrategy struct {
}

// IssueURLStrategy take what the user provided, assume it's a URL take the
// issue number from the end
type IssueURLStrategy struct {
}

// GetIssue extracts the issue number from the issue
func (c IssueURLStrategy) GetIssue(rawIssue string) (string, error) {
	issueURL, err := url.Parse(rawIssue)
	if err != nil {
		return "", errors.Wrap(err, "issue url invalid")
	} else if issueURL.Host == "" {
		return "", errors.New("no host provided, not a valid url")
	}

	splitPath := strings.Split(issueURL.Path, "/")

	return splitPath[len(splitPath)-1], nil
}

// GetIssue extracts the issue number from the issue
func (c IssueLiteralStrategy) GetIssue(rawIssue string) (string, error) {
	return rawIssue, nil
}

// MakeIssueStrategy make an issue strategy based on the URL provided
func MakeIssueStrategy(
	issueURL *url.URL,
) IssueStrategy {
	if issueURL != nil && issueURL.Host != "" {
		return IssueURLStrategy{}
	}

	return IssueLiteralStrategy{}
}
