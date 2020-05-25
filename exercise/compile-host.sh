#!/usr/bin/env bash
set -o errexit

main () {
  cd "$(dirname "$0")"
  cd ..

  go build -o exercise/.bin/terraform-provider-rackcorp
}

main "$@"
