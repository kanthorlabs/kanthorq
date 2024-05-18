#!/bin/bash
set -e

DESTINATION=${DESTINATION:-"/tmp/kanthorlabs/common"}

rm -rf $DESTINATION
git clone git@github.com:kanthorlabs/common.git $DESTINATION

# move to destination
pushd $DESTINATION
COMMIT=$(git rev-parse --short HEAD)
# move back to original directory
popd

# upgrade
echo "upgrading"
go get -u github.com/kanthorlabs/common@$COMMIT && go mod vendor

echo "cleaning"
rm -rf $DESTINATION