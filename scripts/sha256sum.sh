#!/bin/bash
set -e

find . -type f \( -path './entities/*.go' -o -path './core/*.go' -o -path './publisher/*.go' -o -path './subscriber/*.go' \) -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1
