#!/bin/bash
set -e

find . -type f \( -path './core/*.go' -o -path './entities/*.go' -o -path './publisher/*.go' -o -path './puller/*.go' -o -path './subscriber/*.go' -o -path './*.go' \) -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1
