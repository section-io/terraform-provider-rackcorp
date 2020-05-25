FROM golang:1.13.7 as build

ENV CGO_ENABLED=0

WORKDIR /src/terraform-provider-rackcorp

COPY go.mod go.sum tools.go ./
RUN go mod download

# explicitly install tools to improve Docker re-build times
RUN grep '^[[:space:]]*_[[:space:]]\+"[^"]\+"' tools.go | cut -d'"' -f2 | xargs -t go install

COPY . .

RUN gofmt -e -s -d . 2>&1 | tee /gofmt.out && test ! -s /gofmt.out
RUN go vet ./...
RUN golint -set_exit_status ./...

RUN errcheck ./...

RUN go install ./...

RUN go test -v ./...

### END FROM build

FROM hashicorp/terraform:0.11.10 as base

RUN printf 'providers {\n  rackcorp = "/go/bin/terraform-provider-rackcorp"\n}\n' >/root/.terraformrc

# TODO add version to file name as per:
#  https://www.terraform.io/docs/configuration/providers.html#plugin-names-and-versions
COPY --from=build /go/bin/terraform-provider-rackcorp /go/bin/

FROM base as test

WORKDIR /work
COPY example.tf ./main.tf

RUN terraform fmt -diff -check ./

# https://www.terraform.io/docs/internals/debugging.html
ARG TF_LOG=WARN

RUN terraform init && \
  terraform plan -out=a.tfplan

FROM base
