#!/usr/bin/env bash

# AVOID INVOKING THIS SCRIPT DIRECTLY -- USE `make test`

set -euxo pipefail

GO111MODULE=on \
    go test -v \
    ./devices/... \
    ./examples/... \
    ./features/... \
    ./protocols/... \
    ./version/...
