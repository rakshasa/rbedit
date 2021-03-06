#!/bin/bash

set -eux

BUILD_IMAGE="${BUILD_IMAGE:-build-env}"
BUILD_MARKDOWN="${BUILD_MARKDOWN:-no}"
BUILD_DOCS="${BUILD_DOCS:-no}"
BUILD_DIR="${BUILD_DIR:-./build}"

TARGET_OS="${TARGET_OS:-}"
TARGET_ARCH="${TARGET_ARCH:-amd64}"

case "${TARGET_OS:-}" in
  darwin|linux|windows)
    echo "TARGET_OS=${TARGET_OS}"
    ;;
  "")
    if [[ ${OSTYPE} =~ ^darwin ]]; then
      TARGET_OS=darwin
    elif [[ ${OSTYPE} =~ ^linux-gnu ]]; then
      TARGET_OS=linux
    else
      echo "could not detect a valid target OS"
      exit 1
    fi
    ;;
  *)
    echo "unknown target OS: ${TARGET_OS}"
    exit 1
    ;;
esac

if [[ "${BUILD_DOCS}" == "yes" ]]; then
  BUILD_MARKDOWN="yes"
elif [[ "${BUILD_DOCS}" != "no" ]]; then
  echo "BUILD_DOCS must be either 'yes' or 'no'"
  exit 1
fi

if ! [[ "${BUILD_MARKDOWN}" =~ ^(yes)|(no)$ ]]; then
  echo "BUILD_MARKDOWN must be either 'yes' or 'no'"
  exit 1
fi

project_root="$(cd "$(cd "$( dirname "${BASH_SOURCE[0]}" )" && git rev-parse --show-toplevel)" >/dev/null 2>&1 && pwd)"
readonly project_root

readonly container="rtdo-build-rbedit"
readonly rbedit_image="rtdo/rbedit"

cleanup() {
  local -r retval="$?"
  set +eu

  docker rm -f "${container}"
  rm -rf "${build_dir}"

  set +x

  if (( retval == 0 )); then
    echo
    echo "***********************"
    echo "*** Build Succeeded ***"
    echo "***********************"
    echo
  else
    echo
    echo "********************"
    echo "*** Build Failed ***"
    echo "********************"
    echo
  fi

  exit "${retval}"
}
trap cleanup EXIT

dockerfile_no_builder() {
  sed -n -e '/ AS rbedit-builder$/,$p' dockerfile | sed "s|^FROM build-env AS rbedit-builder\$|FROM \"${BUILD_IMAGE}\" AS rbedit-builder|"
}

build_dir=$(mktemp -d); readonly build_dir

( cd "${build_dir}"

  git clone --depth 1 file://"${project_root}" ./

  if [[ "${BUILD_IMAGE}" == "build-env" ]]; then
    readonly build_file="./dockerfile"
  else
    echo "Using '${BUILD_IMAGE}' as the build image."

    readonly build_file="./dockerfile.build"
    dockerfile_no_builder > "${build_file}"

    echo
    cat "${build_file}"
    echo
  fi

  docker build \
    --progress plain \
    --tag "${rbedit_image}"\
    --target "rbedit-builder" \
    --file "${build_file}" \
    --build-arg "TARGET_OS=${TARGET_OS}" \
    --build-arg "TARGET_ARCH=${TARGET_ARCH}" \
    --build-arg "BUILD_MARKDOWN=${BUILD_MARKDOWN}" \
    .
)

( cd "${project_root}"

  docker create -i --rm \
    --name "${container}" \
    "${rbedit_image}"

  mkdir -p "${BUILD_DIR}/"
  docker cp "${container}:/rbedit-${TARGET_OS}-${TARGET_ARCH}" "${BUILD_DIR}/"

  if [[ "${BUILD_MARKDOWN}" == "yes" ]]; then
    docker cp "${container}:/rbedit-markdown-${TARGET_OS}-${TARGET_ARCH}" "${BUILD_DIR}/"
  fi

  if [[ "${BUILD_DOCS}" == "yes" ]]; then
    if ! "${BUILD_DIR}/rbedit-markdown-${TARGET_OS}-${TARGET_ARCH}" &> /dev/null; then
      echo "could not run ${BUILD_DIR}/rbedit-markdown-${TARGET_OS}-${TARGET_ARCH}"
      exit 1
    fi

    rm -rf ./docs/cli
    mkdir -p ./docs/cli

    "${BUILD_DIR}/rbedit-markdown-${TARGET_OS}-${TARGET_ARCH}" ./docs/cli

    git add ./docs/cli
  fi
)

success="yes"
