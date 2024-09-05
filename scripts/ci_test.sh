#!/bin/bash
set -e

export $(grep -v '^#' .env | xargs)

CI=${CI:-""}
SCRIPT_DIR=$(dirname "$0")
CHECKSUM_FILE=./checksum
CHECKSUM_NEW=$($SCRIPT_DIR/sha256sum.sh)
CHECKSUM_OLD=$(cat $CHECKSUM_FILE || true)

if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "--> coverage"
  go test -timeout 1m30s --count=1 -cover -coverprofile cover.out \
    github.com/kanthorlabs/kanthorq \
    github.com/kanthorlabs/kanthorq/entities \
    github.com/kanthorlabs/kanthorq/core \
    github.com/kanthorlabs/kanthorq/publisher \
    github.com/kanthorlabs/kanthorq/puller \
    github.com/kanthorlabs/kanthorq/subscriber
fi

if [ "$CI" = "" ]; then
  echo -n $CHECKSUM_NEW > ./checksum
fi
