#!/bin/bash

set -eux

project_root="$(cd "$(cd "$( dirname "${BASH_SOURCE[0]}" )" && git rev-parse --show-toplevel)" >/dev/null 2>&1 && pwd)"
readonly project_root

readonly container="rbedit-update-vendor"
readonly build_image="rtdo/build/rbedit"

readonly dependencies=(
  # github.com/jackpal/bencode-go@v1.0.0
  github.com/rakshasa/bencode-go@v1.0.2
  github.com/spf13/cobra@v1.2.1
)

build_dir=$(mktemp -d); readonly build_dir

cleanup() {
  local -r retval="$?"
  set +eu

  docker rm "${container}"
  docker rmi "${build_image}"

  rm -rf "${build_dir}"

  set +x

  if [[ "${success:-no}" == "yes" ]]; then
    echo
    echo "*******************************"
    echo "*** Vendor Update Succeeded ***"
    echo "*******************************"
    echo
  else
    echo
    echo "****************************"
    echo "*** Vendor Update Failed ***"
    echo "****************************"
    echo
  fi

  exit "${retval}"
}
trap cleanup EXIT 1 3 6 8 11 14 15 20 26

cd "${project_root}"

git clone --depth 1 file://"${project_root}" "${build_dir}"

docker build \
  --progress plain \
  --file "./Dockerfile" \
  --target "build-env" \
  --tag "${build_image}"\
  .

docker run -i \
  --name "${container}" \
  --volume "${build_dir}":/build/go/src/github.com/rakshasa/rbedit \
  "${build_image}" \
  /bin/bash - <<"EOF"
#!/bin/bash
set -euxo pipefail

cd /build/go/src/github.com/rakshasa/rbedit/

rm -rf ./go.{mod,sum} ./vendor/

go clean -cache
go mod init github.com/rakshasa/rbedit

for dep in "${dependencies[@]}"; do
  go get -u -v "${dep}"
done

go mod tidy -v
go mod vendor -v

set +x
echo
echo "+----------------------+"
echo "| Vendor Files Created |"
echo "+----------------------+"
echo

EOF

ls "${build_dir}"

rm -rf ./{go.mod,go.sum,vendor}
cp -r "${build_dir}"/{go.mod,go.sum,vendor} ./

success="yes"
