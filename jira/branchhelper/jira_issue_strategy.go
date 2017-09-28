package branchhelper

import (
	"errors"
	"net/url"
	"strings"
)

type IssueStrategy interface {
	GetIssue(rawIssue string) (string, error)
}

type IssueLiteralStrategy struct {
}

type IssueUrlStrategy struct {
}

func (c IssueUrlStrategy) GetIssue(rawIssue string) (string, error) {
	issueUrl, err := url.Parse(rawIssue)

	if err != nil {
		return "", err
	}
	if issueUrl.Host == "" {
		return "", errors.New("no host provided, not a valid url")
	}

	splitPath := strings.Split(issueUrl.Path, "/")

	return splitPath[len(splitPath)-1], nil
}

func (c IssueLiteralStrategy) GetIssue(rawIssue string) (string, error) {
	return rawIssue, nil
}
