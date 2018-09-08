#!/usr/bin/env bash

# AVOID INVOKING THIS SCRIPT DIRECTLY -- USE `make dep`

# This script resolves dependencies by building with support for Go modules
# enabled. Binaries are subsequently deleted since this script does not bother
# with setting build flags correctly.

set -euxo pipefail

export GO111MODULE=on

go mod vendor
go mod tidy
