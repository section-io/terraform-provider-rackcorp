FROM golang:1.12 as build

ENV CGO_ENABLED=0

WORKDIR /go/src/app

# explicitly install dependencies to improve Docker re-build times
RUN go get -v \
  github.com/kisielk/errcheck \
  golang.org/x/lint/golint \
  gopkg.in/h2non/gock.v1

# Use specific version of terraform
RUN mkdir -p "${GOPATH}/src/github.com/hashicorp/" && \
  cd "${GOPATH}/src/github.com/hashicorp/" && \
  git clone --verbose --branch v0.12.5 --depth 1 https://github.com/hashicorp/terraform

WORKDIR /app

# BEGIN pre-install dependencies to reduce time for each code change to build
COPY go.mod go.sum ./
RUN go mod download
# END

COPY main.go .
COPY rackcorp ./rackcorp/

RUN gofmt -e -s -d . 2>&1 | tee /gofmt.out && test ! -s /gofmt.out
RUN go vet .
RUN golint -set_exit_status ./...
RUN errcheck ./...

RUN go install ./...

RUN go test -v ./...

### END FROM build

FROM hashicorp/terraform:0.12.5

RUN mkdir -p /work/ && \
  printf 'providers {\n  rackcorp = "/go/bin/terraform-provider-rackcorp"\n}\n' >/root/.terraformrc

WORKDIR /work

COPY --from=build /go/bin/terraform-provider-rackcorp /go/bin/

COPY example.tf ./main.tf

RUN terraform fmt -diff -check ./

# https://www.terraform.io/docs/internals/debugging.html
ARG TF_LOG=WARN

RUN terraform init && \
  terraform plan -out=a.tfplan
