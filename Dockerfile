FROM golang:1.9

ENV CGO_ENABLED=0

WORKDIR /go/src/app

RUN go get -u github.com/kisielk/errcheck

# explicitly install dependencies to improve Docker re-build times
RUN go get -v -d \
  github.com/hashicorp/terraform

RUN mkdir -p /go/src/github.com/section-io/ && \
  ln -s /go/src/app /go/src/github.com/section-io/terraform-provider-rackcorp

COPY *.go /go/src/app/
COPY rackcorp /go/src/app/rackcorp

RUN gofmt -e -s -d /go/src/app 2>&1 | tee /gofmt.out && test ! -s /gofmt.out
RUN go tool vet /go/src/app
RUN errcheck ./...

RUN go-wrapper install

RUN cd "/go/src/$(go list -e -f '{{.ImportComment}}')" && \
  go test -bench=. -v ./...
