#!/usr/bin/env bash

# AVOID INVOKING THIS SCRIPT DIRECTLY -- USE `make test`

set -euxo pipefail

GO111MODULE=on \
    go test -timeout 30s -race -coverprofile=coverage.txt -covermode=atomic \
    ./examples/... \
    ./features/... \
    ./products/... \
    ./protocols/...
