#!/usr/bin/env bash
set -o errexit

use_local_credentials () {
  while IFS='=' read -r key value; do
    case "${key}" in
      APIUUID)
        export RACKCORP_API_UUID="${value}"
        ;;
      APISECRET)
        export RACKCORP_API_SECRET="${value}"
        ;;
    esac
  done <"${HOME}/.config/rackcorp/apikey"

  if [[ -z "${RACKCORP_CUSTOMER_ID}" ]]; then
    printf 'Error: RACKCORP_CUSTOMER_ID environment variable required.\n' >&2
    return 1
  fi
}

terraform () {
  local version=0.11.10
  local url="https://releases.hashicorp.com/terraform/${version}/terraform_${version}_linux_amd64.zip"
  local binary="${SCRIPT_ROOT}/.bin/terraform-${version}"

  if ! "${binary}" version >/dev/null; then
    local workdir=$(mktemp -d)
    wget "${url}" -O "${workdir}/zip"
    unzip "${workdir}/zip" -d "${workdir}"
    mkdir -p "$(dirname "${binary}")"
    mv "${workdir}/terraform" "${binary}"
    rm "${workdir}/zip"
    rm -d "${workdir}"
  fi

  "${binary}" "$@"
}

main () {
  cd "$(dirname "$0")"
  declare -g SCRIPT_ROOT="${PWD}"

  use_local_credentials

  terraform version

  for scenario in custom-image ubuntu-image ubuntu-install; do

    cd "${SCRIPT_ROOT}/${scenario}"

    terraform fmt
    terraform init

    terraform plan -out .tfplan

    terraform apply .tfplan

    sleep 10

    terraform plan -destroy -out .tfplan

    terraform apply .tfplan

  done
}

main "$@"
