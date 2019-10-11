#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

export GO111MODULE=on

# make sure we're in the correct directory
cd "$(dirname "$0")"

APP="jokesontap"

if [[ ! -d ./vendor ]];then
  echo 'vendoring required: use `go mod vendor` to push dependencies to vendor directory'
fi

# assign version from VERSION file at project root
VERSION=$(head -n 1 VERSION)
if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+.*$ ]]; then
  echo "version '$VERSION' does not match expected format, inspect VERSION file"
  exit 1
fi

# test
if ! go test -v ./... -mod=vendor; then
  echo "go test failed, build process aborted"
fi

# clean up existing binaries
rm -rf ./bin

# build
if ! go build -ldflags "-X main.buildVersion=${VERSION}" -v -mod=vendor -o "./bin/${APP}" ./cmd/srv/; then
  exit 1
fi

echo "binary written to ./bin"
