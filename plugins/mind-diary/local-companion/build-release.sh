#!/bin/sh

set -eu

if [ "$#" -ne 2 ]; then
  printf '%s\n' 'usage: build-release.sh <immutable-plugin-version> <output-directory>' >&2
  exit 2
fi

plugin_version=$1
output_directory=$2

if ! printf '%s\n' "$plugin_version" | LC_ALL=C grep -Eq '^0\.1\.0\+codex\.[0-9]{14}$'; then
  printf '%s\n' 'the plugin version must be an immutable 0.1.0+codex.YYYYMMDDhhmmss value' >&2
  exit 2
fi

script_directory=$(CDPATH= cd -- "${0%/*}" && pwd)
mkdir -p "$output_directory"
output_directory=$(CDPATH= cd -- "$output_directory" && pwd)

build_architecture() {
  architecture=$1
  output_name=$2
  (
    cd "$script_directory"
    CGO_ENABLED=0 GOOS=darwin GOARCH="$architecture" GOFLAGS= \
      go build -trimpath -buildvcs=false \
      -ldflags="-s -w -X=main.buildVersion=$plugin_version" \
      -o "$output_directory/$output_name" .
  )
}

build_architecture arm64 mind-diary-local-darwin-arm64
build_architecture amd64 mind-diary-local-darwin-amd64
