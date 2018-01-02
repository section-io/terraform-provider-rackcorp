FROM golang:1.9 as build

ENV CGO_ENABLED=0

WORKDIR /go/src/app

RUN go get -u github.com/kisielk/errcheck

# explicitly install dependencies to improve Docker re-build times
RUN go get -v \
  github.com/hashicorp/terraform \
  github.com/pkg/errors \
  github.com/stretchr/testify/assert \
  gopkg.in/h2non/gock.v1

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

### END FROM build

FROM hashicorp/terraform:0.11.1

RUN mkdir -p /work/ && \
  printf 'providers {\n  rackcorp = "/go/bin/terraform-provider-rackcorp"\n}\n' >/root/.terraformrc

WORKDIR /work

COPY --from=build /go/bin/terraform-provider-rackcorp /go/bin/

COPY example.tf ./main.tf

# https://www.terraform.io/docs/internals/debugging.html
ARG TF_LOG=WARN

RUN terraform init && \
  terraform plan -out=a.tfplan && \
  terraform apply a.tfplan && \
  terraform plan -out=b.tfplan && \
  terraform plan -destroy -out=destroy.tfplan && \
  terraform apply destroy.tfplan
