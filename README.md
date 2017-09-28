# jira-branch-helper

[![Go Report Card](https://goreportcard.com/badge/github.com/PurpleBooth/jira-branch-helper)][3]
[![codebeat badge](https://codebeat.co/badges/60d295f5-72cf-42cf-9fd8-db8fab7389ac)][4]

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
  Build a string that can be used for a branch name from the details in a Jira ticket

  Example usage

  $ jira-branch-helper https://example.com/jira/browse/TST-123
  tst-123-ticket-title-goes-here

  $ jira-branch-helper TST-123
  tst-123-ticket-title-goes-here

  Environment variables may be used in place of flags, parameters, see parameters with [$ENV_NAME_HERE] at the end.

  The following functions are available for templating

  * "Trim" - Remove whitespace at the start and end of the function (no parameters)
  * "ToLower" - Lower case your string (no parameters)
  * "ToUpper" - Upper case your string (no parameters)
  * "Replace" - Replace characters e.g. {{.Fields.Summary | Replace "A" "B" }}, would replace A with B
  * "KebabCase" - e.g. developments-phase-1-implement-feature--bume (no parameters)
  * "LowerSnakeCase" - e.g. developments_phase_1_implement_feature__bume (no parameters)
  * "LowerCamelCase" - e.g. developmentsPhase1ImplementFeatureBume (no parameters)
  * "ScreamingKebabCase" - e.g. DEVELOPMENTS-PHASE-1-IMPLEMENT-FEATURE--BUME (no parameters)
  * "ScreamingSnakeCase" - e.g. DEVELOPMENTS_PHASE_1_IMPLEMENT_FEATURE__BUME (no parameters)
  * "UpperCamelCase" - e.g. DevelopmentsPhase1ImplementFeatureBume (no parameters)
  * "UpperKebabCase" - e.g. Developments_Phase_1_Implement_Feature__Bume (no parameters)
  
  The template format is as described here https://golang.org/pkg/text/template/

USAGE:
   jira-branch-helper [global options] command [command options] [ISSUE-NUMBER OR ISSUE-URL]

VERSION:
   v0.1.0

AUTHOR:
   Billie Alice Thompson <billie@purplebooth.co.uk>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --jira-username value  The username to authenticate as on Jira [$JIRA_BRANCH_HELPER_USERNAME]
   --jira-password value  The password to authenticate as on Jira [$JIRA_BRANCH_HELPER_PASSWORD]
   --jira-endpoint value  Jira's URL [$JIRA_BRANCH_HELPER_ENDPOINT]
   --template value       The template to use to generate the branch name (default: "{{.Key | ToLower }}-{{.Fields.Summary | Trim | KebabCase }}") [$JIRA_BRANCH_HELPER_TEMPLATE]
   --help, -h             show help
   --version, -v          print the version

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

[2]: https://godoc.org/github.com/PurpleBooth/jira-branch-helper/jira-branch-helper
[3]: https://goreportcard.com/report/github.com/PurpleBooth/jira-branch-helper
[4]: https://codebeat.co/projects/github-com-purplebooth-jira-branch-helper-master