#!/bin/bash

set -eu

if (( ${#} != 2 )); then
  echo "required arguments: RBEDIT_PATH=./build/rbedit COUNT=1000 $(basename "${0}") SRC-TORRENT DEST-DIR" 1>&2
  exit 1
fi

COUNT="${COUNT:-1000}"
RBEDIT_PATH="${RBEDIT_PATH:-./build/rbedit}"

readonly src_torrent="${1}"
readonly dest_dir="${2}"

echo "Generating torrent test files"
echo
echo "RBEDIT_PATH: ${RBEDIT_PATH}"
echo "COUNT: ${COUNT}"
echo "SRC-TORRENT: ${src_torrent}"
echo "DEST-DIR: ${dest_dir}"

if ! [[ -f "${src_torrent}" ]]; then
  echo "torrent source is not a file" 1>&2
  exit 1
fi

if ! [[ -d "${dest_dir}" ]]; then
  echo "destination is not a directory" 1>&2
  exit 1
fi

if ! [[ "${COUNT}" =~ ^[1-9][0-9]+$ ]]; then
  echo "COUNT is not a positive number" 1>&2
  exit 1
fi

if ! [[ -x "${RBEDIT_PATH}" ]]; then
  echo "rbedit not found or executable" 1>&2
  exit 1
fi

prefix_depth=0

for (( count=COUNT; count > 1024; count=count/1024 )); do
  : $(( prefix_depth++ ))
done

echo "PREFIX-DEPTH: ${prefix_depth}"
echo
echo -n "generating"

for (( c=0; c < COUNT; c++ )); do
  if (( c % 100 == 0 )); then
    echo -n "."
  fi

  tmp_torrent="$(mktemp)"

  cp "${src_torrent}" "${tmp_torrent}"

  "${RBEDIT_PATH}" put --input "${tmp_torrent}" --inplace --int "${c}" info test-value

  info_hash="$("${RBEDIT_PATH}" hash info --input "${tmp_torrent}")"
  dest_prefix="${dest_dir}"

  for (( d=0; d < prefix_depth; d++ )); do
    dest_prefix+="/${info_hash:${d}:1}"
  done

  if ! [[ -d "${dest_prefix}" ]]; then
    mkdir -p "${dest_prefix}"
  fi

  mv "${tmp_torrent}" "${dest_prefix}/${info_hash}"
done

echo
echo
echo "Finished generating ${COUNT} torrents"
