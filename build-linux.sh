#!/usr/bin/env bash

cd "$(dirname "$0")" || exit

image_name=terraform-provider-rackcorp-builder:latest

rm -rf bin || true
mkdir bin || exit

docker build -t "${image_name}" . || exit

cid=$(docker create "${image_name}") || exit

docker cp "${cid}:/root/.terraform.d/plugins/" -| tar --extract --file=- --strip-components=1 --directory=./bin || exit

docker rm "${cid}" >/dev/null
