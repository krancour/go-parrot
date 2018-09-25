#!/usr/bin/env bash

# AVOID INVOKING THIS SCRIPT DIRECTLY -- USE `make test`

set -euxo pipefail

export GO111MODULE=off

go test -v \
    ./devices/... \
    ./examples/... \
    ./features/... \
    ./protocols/... \
    ./version/...
