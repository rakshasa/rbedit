#!/bin/bash

set -eux

project_root="$(cd "$(cd "$( dirname "${BASH_SOURCE[0]}" )" && git rev-parse --show-toplevel)" >/dev/null 2>&1 && pwd)"
readonly project_root

readonly container="rtdo-build-rbedit"
readonly rbedit_image="rtdo/rbedit"

cleanup() {
  local -r retval="$?"
  set +eux

  if (( retval == 0 )); then
    echo
    echo "*************************"
    echo "*** Release Succeeded ***"
    echo "*************************"
    echo
  else
    echo
    echo "**********************"
    echo "*** Release Failed ***"
    echo "**********************"
    echo
  fi

  exit "${retval}"
}
trap cleanup EXIT

build_dir="$(mktemp -d)"

( cd "${project_root}"

  BUILD_DOCS=yes ./scripts/build.sh

  BUILD_DIR="${build_dir}" TARGET_OS=linux ./scripts/build.sh
  BUILD_DIR="${build_dir}" TARGET_OS=darwin ./scripts/build.sh
  BUILD_DIR="${build_dir}" TARGET_OS=windows ./scripts/build.sh
)

( cd "${build_dir}"; ls .

  zip rbedit-darwin-amd64.zip rbedit-darwin-amd64
  zip rbedit-linux-amd64.zip rbedit-linux-amd64
  zip rbedit-windows-amd64.zip rbedit-windows-amd64
)

set +x

printf "\n%s" "${build_dir}/"*.zip
printf "\n"
