#!/bin/bash

set -eux

: "${RBEDIT_ARCH:-linux}"
: "${RBEDIT_BUILD_IMAGE:-build-env}"

project_root="$(cd "$(cd "$( dirname "${BASH_SOURCE[0]}" )" && git rev-parse --show-toplevel)" >/dev/null 2>&1 && pwd)"
readonly project_root

readonly container="rtdo-build-rbedit"
readonly rbedit_image="rtdo/rbedit"

build_dir=$(mktemp -d); readonly build_dir

case "${RBEDIT_ARCH}" in
  "darwin")
    ;;
  "linux")
    ;;
  *)
    echo "invalid RBEDIT_ARCH value: '${RBEDIT_ARCH}'"
    exit 1
    ;;
esac

cleanup() {
  local -r retval="$?"
  set +eu

  docker rm -f "${container}"

  rm -rf "${build_dir}"

  set +x

  if [[ "${success:-no}" == "yes" ]]; then
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
trap cleanup EXIT 1 3 6 8 11 14 15 20 26

dockerfile_no_builder() {
  sed -n -e '/ AS rbedit-builder$/,$p' Dockerfile | sed "s|^FROM build-env AS rbedit-builder\$|FROM \"${RBEDIT_BUILD_IMAGE}\" AS rbedit-builder|"
}

cd "${build_dir}"

git clone --depth 1 file://"${project_root}" ./

if [[ "${RBEDIT_BUILD_IMAGE}" == "build-env" ]]; then
  docker build \
    --progress plain \
    --file "./Dockerfile" \
    --target "rbedit" \
    --tag "${rbedit_image}"\
    --build-arg "TARGET_ARCH=${RBEDIT_ARCH}" \
    .

  readonly build_file="./Dockerfile"
else
  echo "Using '${RBEDIT_BUILD_IMAGE} build image."

  readonly build_file="./dockerfile.build"
  dockerfile_no_builder > "${build_file}"

  echo
  cat "${build_file}"
  echo
fi

docker build \
  --progress plain \
  --file "${build_file}" \
  --target "rbedit" \
  --tag "${rbedit_image}"\
  --build-arg "TARGET_ARCH=${RBEDIT_ARCH}" \
  .

docker create -i --rm \
  --name "${container}" \
  "${rbedit_image}"

mkdir -p "${project_root}/build/"
docker cp "${container}:/rbedit" "${project_root}/build/"

success="yes"
