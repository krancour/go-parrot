#!/usr/bin/env bash

# AVOID INVOKING THIS SCRIPT DIRECTLY -- USE `make lint`

set -euxo pipefail

# gometalinter doesn't seem to work with modules yet, so we will vendor the
# modules for the sake of linting. These are NOT tracked in source control
# since it is not customary for libraries to vendor dependencies.

GO111MODULE=on go mod vendor

GO111MODULE=off \
  golangci-lint run \
  ./controllers/... \
	./examples/... \
	./features/... \
  ./protocols/...
