#!/usr/bin/env bash
readonly SCRIPT_ROOT=$(pwd)
terraform_image=gcr.io/section-io/terraform-provider-rackcorp:latest
# ensure latest
docker pull "${terraform_image}"

terraform () {
  docker run \
    --rm \
    -v "${SCRIPT_ROOT}/:/app/:rw" \
    -e RACKCORP_API_UUID="$RACKCORP_API_UUID" \
    -e RACKCORP_API_SECRET="$RACKCORP_API_SECRET" \
    -e RACKCORP_CUSTOMER_ID="$RACKCORP_CUSTOMER_ID" \
    -w /app/ \
    "${terraform_image}" \
      "$@"
}

terraform fmt -diff -write=false

terraform init -input=false

terraform plan \
  -out apply.tfplan

if [ 'true' = "${SKIP_APPLY}" ]
then
  printf 'Skipping "terraform apply".\n'
  exit
fi

terraform apply apply.tfplan

