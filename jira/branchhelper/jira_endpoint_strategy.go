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
)

// EndpointStrategy Get an endpoint URL from an issue URL
type EndpointStrategy interface {
	GetEndpoint(url *url.URL) *url.URL
}

// EndpointCombinedStrategy is used for instances Jira is deployed with other
// products, so has the "/jira" prefix e.g. "/jira/browse/TST-1"
type EndpointCombinedStrategy struct {
}

func (c EndpointCombinedStrategy) applicable(issueURL *url.URL) bool {
	return strings.HasPrefix(issueURL.Path, "/jira/browse/")
}

// GetEndpoint Try to get a endpoint URL from a issue url using the
// EndpointCombinedStrategy
func (c EndpointCombinedStrategy) GetEndpoint(issueURL *url.URL) *url.URL {
	if !c.applicable(issueURL) {
		return nil
	}

	return makeURL(issueURL, "/jira")
}

// EndpointSoloStrategy is used for instances Jira is deployed alone, so has no
// prefix to issues e.g. "/browse/TST-1"
type EndpointSoloStrategy struct {
}

func (c EndpointSoloStrategy) applicable(issueURL *url.URL) bool {
	return strings.HasPrefix(issueURL.Path, "/browse/")
}

// GetEndpoint Try to get a endpoint URL from a issue url using the
// EndpointSoloStrategy
func (c EndpointSoloStrategy) GetEndpoint(issueURL *url.URL) *url.URL {
	if !c.applicable(issueURL) {
		return nil
	}

	return makeURL(issueURL, "/")
}

func makeURL(sourceURL *url.URL, path string) *url.URL {

	return &url.URL{
		Path:   path,
		Host:   sourceURL.Host,
		Scheme: sourceURL.Scheme,
	}
}

// GuessEndpointURL Guesses the endpoint URL from a issue URL
func GuessEndpointURL(issueURL *url.URL) string {
	if issueURL == nil {
		return ""
	}

	urlHelpers := []EndpointStrategy{
		EndpointCombinedStrategy{},
		EndpointSoloStrategy{},
	}

	for i := range urlHelpers {
		possibleEndpointURL := urlHelpers[i].GetEndpoint(issueURL)

		if possibleEndpointURL != nil {
			return possibleEndpointURL.String()
		}
	}

	return ""
}
