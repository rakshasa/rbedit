#!/bin/bash

set -eux

project_root="$(cd "$(cd "$( dirname "${BASH_SOURCE[0]}" )" && git rev-parse --show-toplevel)" >/dev/null 2>&1 && pwd)"
readonly project_root

readonly rbedit_markdown_image="rtdo/rbedit-markdown"

build_dir=$(mktemp -d); readonly build_dir
wiki_dir=$(mktemp -d); readonly wiki_dir

cleanup() {
  local -r retval="$?"
  set +eu

  rm -rf "${build_dir}" "${wiki_dir}"

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

( cd "${build_dir}"

  git clone --depth 1 file://"${project_root}" ./

  docker build \
    --tag "${rbedit_markdown_image}"\
    --progress plain \
    --file ./dockerfile \
    --target "rbedit-markdown" \
    --build-arg "BUILD_MARKDOWN=yes" \
    .
)

( cd "${project_root}"

  if ! docker run -i --rm "${rbedit_markdown_image}" &> /dev/null; then
    echo "could not run: docker run -i --rm \"${rbedit_markdown_image}\""
    exit 1
  fi

  rm -rf ./docs/cli
  mkdir -p ./docs/cli

  docker run \
    --interactive \
    --rm \
    --volume "${project_root}/docs/cli:/docs/cli" \
    "${rbedit_markdown_image}" \
    /docs/cli

  git add ./docs/cli
)

( cd "${wiki_dir}"

  git clone git@github.com:rakshasa/rbedit.wiki.git ./
  git rm ./rbedit*.md || :

  cp "${project_root}"/docs/cli/rbedit*.md ./
  sed -i '' -e 's/\[\([a-z -]*\)\](\([a-z _-]*\).md)/[[\1\|\2]]/' ./rbedit*.md

  git add ./rbedit*.md

  if git commit -m "Updated cobra-generated documents."; then
    git push
  fi
)
