#!/usr/bin/env bash
set -o errexit

main () {
  cd "$(dirname "$0")"

  image_name=section.io.invalid/section-io/terraform-provider-rackcorp

  docker build -t "${image_name}" ./..

  cid=$(docker create "${image_name}")

  docker cp "${cid}:/go/bin/terraform-provider-rackcorp" ./.bin/

  docker rm -f "${cid}" >/dev/null
}

main "$@"
