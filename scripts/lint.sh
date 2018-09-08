#!/usr/bin/env bash

# AVOID INVOKING THIS SCRIPT DIRECTLY -- USE `make lint`

set -euxo pipefail

export GO111MODULE=off

gometalinter cmd/... pkg/... \
	--concurrency=1 \
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
