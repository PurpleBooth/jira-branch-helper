package branchhelper

import (
	"net/url"
	"strings"
)

type EndpointStrategy interface {
	GetEndpoint(url *url.URL) *url.URL
}

type EndpointCombinedStrategy struct {
}

func (c EndpointCombinedStrategy) applicable(issueUrl *url.URL) bool {
	return strings.HasPrefix(issueUrl.Path, "/jira/browse/")
}
func (c EndpointCombinedStrategy) GetEndpoint(issueUrl *url.URL) *url.URL {
	if !c.applicable(issueUrl) {
		return nil
	}

	return makeUrl(issueUrl, "/jira")
}

type EndpointSoloStrategy struct {
}

func (c EndpointSoloStrategy) applicable(issueUrl *url.URL) bool {
	return strings.HasPrefix(issueUrl.Path, "/browse/")
}

func (c EndpointSoloStrategy) GetEndpoint(issueUrl *url.URL) *url.URL {
	if !c.applicable(issueUrl) {
		return nil
	}

	return makeUrl(issueUrl, "/")
}

func makeUrl(sourceUrl *url.URL, path string) *url.URL {

	return &url.URL{
		Path:   path,
		Host:   sourceUrl.Host,
		Scheme: sourceUrl.Scheme,
	}
}
