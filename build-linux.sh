#!/usr/bin/env bash

cd "$(dirname "$0")" || exit

image_name=terraform-provider-rackcorp-builder:latest

mkdir bin || exit

docker build -t "${image_name}" . || exit

cid=$(docker create "${image_name}") || exit

docker cp "${cid}:/go/bin/terraform-provider-rackcorp" ./bin/ || exit

docker rm "${cid}" >/dev/null
