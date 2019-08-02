FROM golang:1.12 as build

ENV CGO_ENABLED=0

WORKDIR /go/src/app

# explicitly install dependencies to improve Docker re-build times
RUN go get -v \
  github.com/kisielk/errcheck \
  golang.org/x/lint/golint \
  gopkg.in/h2non/gock.v1

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

FROM hashicorp/terraform:0.12.6

RUN mkdir -p /work/.terraform/plugins /root/.terraform.d/plugins

WORKDIR /work

COPY --from=build /go/bin/terraform-provider-rackcorp /root/.terraform.d/plugins

COPY example.tf ./main.tf

RUN terraform fmt -diff -check ./

# https://www.terraform.io/docs/internals/debugging.html
ARG TF_LOG=WARN

RUN terraform init && \
  terraform plan -out=a.tfplan
