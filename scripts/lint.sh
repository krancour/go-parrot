#!/usr/bin/env bash

# AVOID INVOKING THIS SCRIPT DIRECTLY -- USE `make lint`

set -euxo pipefail

# gometalinter doesn't seem to work with modules yet, so we will vendor the
# modules for the sake of linting. These are NOT tracked in source control
# since it is not customary for libraries to vendor dependencies.

GO111MODULE=on go mod vendor

GO111MODULE=off \
	gometalinter ./devices/... \
	./examples/... \
	./features/... \
	./protocols/... \
	--disable-all \
	--enable gofmt \
	--enable vet \
	--enable vetshadow \
	--enable gotype \
	--enable deadcode \
	--enable golint \
	--enable varcheck \
	--enable structcheck \
	--enable errcheck \
	--enable megacheck \
	--enable ineffassign \
	--enable interfacer \
	--enable unconvert \
	--enable goconst \
	--enable goimports \
	--enable misspell \
	--enable unparam \
	--enable lll \
	--line-length 80 \
	--deadline 240s
