#!/bin/sh
set -e

export $(grep -v '^#' .env | xargs)

CI=${CI:-""}
CHECKSUM_FILE=./checksum
CHECKSUM_NEW=$(find . -type f \( -name '*.go' \) ! -path "./vendor/*" -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $CHECKSUM_FILE || true)

if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "--> coverage"
  go test -timeout 1m30s --count=1 -cover -coverprofile cover.out $(go list ./... | grep github.com/kanthorlabs/kanthorq | grep -v 'github.com/kanthorlabs/kanthorq/\(cmd\|testify\)')
fi

if [ "$CI" = "" ];
then
  find . -type f \( -name '*.go' \) ! -path "./vendor/*" -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > ./checksum
fi
