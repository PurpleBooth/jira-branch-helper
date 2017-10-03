# jira-branch-helper

[![Go Report Card](https://goreportcard.com/badge/github.com/PurpleBooth/jira-branch-helper)][3]
[![codebeat badge](https://codebeat.co/badges/60d295f5-72cf-42cf-9fd8-db8fab7389ac)][4]
[![Build Status](https://travis-ci.org/PurpleBooth/jira-branch-helper.svg?branch=master)](https://travis-ci.org/PurpleBooth/jira-branch-helper)
[![Docker Build Status](https://img.shields.io/docker/build/purplebooth/jira-branchhelper.svg)](https://hub.docker.com/r/purplebooth/jira-branchhelper/)

Build a string that can be used for a branch name from the details in a Jira ticket
  
## Installing

```bash
go get -u github.com/golang/dep/cmd/dep
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
dep ensure
(cd jira/branchhelper && go test)
(cd cmd/jira-branch-helper/ && go install)
```

## Running

```shell
$ jira-branch-helper help
 NAME:
    jira-branch-helper -
   Build a string that can be used for a branch name from the details in a Jira
   ticket

   Example usage

   $ jira-branch-helper https://example.com/jira/browse/TST-123
   tst-123-ticket-title-goes-here

   $ jira-branch-helper TST-123
   tst-123-ticket-title-goes-here

   Environment variables may be used in place of flags, parameters, see
   parameters with [$ENV_NAME_HERE] at the end.

   The following functions are available for templating

   * "Trim"               - Remove whitespace from start and end
   * "ToLower"            - Lower case your string
   * "ToUpper"            - Upper case your string
   * "Replace"            - Replace characters params: for-search, replace-with
   * "KebabCase"          - Switch the casing-to-kebab
   * "LowerSnakeCase"     - Switch the casing_to_snake
   * "LowerCamelCase"     - Switch the casingToCamel
   * "ScreamingKebabCase" - Switch the CASING-TO-KEBAB
   * "ScreamingSnakeCase" - Switch the CASING_TO_SNAKE
   * "UpperCamelCase"     - Switch the CasingToCamel
   * "UpperKebabCase"     - Switch the Casing-To-Kebab

   The template format is as described here
   https://golang.org/pkg/text/template/

   Templates look like this

   {{.Key | ToLower }}-{{.Fields.Summary | Replace "A" "B" | KebabCase }}


 USAGE:
    jira-branch-helper [global options] command [command options] [ISSUE-NUMBER OR ISSUE-URL]

 AUTHOR(S):
    Billie Alice Thompson <billie@purplebooth.co.uk>

 COMMANDS:
      help, h  Shows a list of commands or help for one command

 GLOBAL OPTIONS:
    --jira-basic-auth-username value  Set a basic auth username on HTTP requests to Jira [$JIRA_BRANCH_HELPER_USERNAME_BASIC_AUTH]
    --jira-basic-auth-password value  Set a basic auth password on HTTP requests to Jira [$JIRA_BRANCH_HELPER_PASSWORD_BASIC_AUTH]
    --jira-username value             The username to authenticate as on Jira [$JIRA_BRANCH_HELPER_USERNAME]
    --jira-password value             The password to authenticate as on Jira [$JIRA_BRANCH_HELPER_PASSWORD]
    --jira-endpoint value             Jira's URL [$JIRA_BRANCH_HELPER_ENDPOINT]
    --template value                  The template to use to generate the branch name (default: "{{.Key | ToLower }}-{{.Fields.Summary | Trim | KebabCase }}") [$JIRA_BRANCH_HELPER_TEMPLATE]
    --help, -h                        show help
    --version, -v                     print the version

 COPYRIGHT:

   jira-branch-helper  Copyright (C) 2017  Billie Alice Thompson
   This program comes with ABSOLUTELY NO WARRANTY;  This is free software,
   and you are welcome to redistribute it under certain conditions; see
   LICENSE.md for additional details.

```


```bash
export JIRA_BRANCH_HELPER_ENDPOINT="https://jira.atlassian.com/"
$ jira-branch-helper TRANS-2457
trans-2457-the-language-picker-in-confluence-cloud-should-be-able-to-show-the-languages
```


## Links

* [Go Docs][2]

[2]: https://godoc.org/github.com/PurpleBooth/jira-branch-helper/jira/branchhelper
[3]: https://goreportcard.com/report/github.com/PurpleBooth/jira-branch-helper
[4]: https://codebeat.co/projects/github-com-purplebooth-jira-branch-helper-master