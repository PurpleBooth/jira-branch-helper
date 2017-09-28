FROM golang:1.9

RUN mkdir -p /go/src/github.com/purplebooth/jira-branch-helper
WORKDIR /go/src/github.com/purplebooth/jira-branch-helper
RUN go get github.com/tools/godep
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
COPY . .
RUN dep ensure
RUN (cd mnemonic && go test)
RUN go build  -o binary -ldflags "-linkmode external -extldflags -static" -a cmd/jira-branch-helper/jira-branch-helper.go

FROM scratch
COPY --from=0 /go/src/github.com/purplebooth/jira-branch-helper/binary /jira-branch-helper
ENTRYPOINT ["/jira-branch-helper"]
