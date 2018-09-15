################################################################################
# Version details                                                              #
################################################################################

GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)

ifeq ($(REL_VERSION),)
	EXAMPLES_VERSION := devel
else
	EXAMPLES_VERSION := $(REL_VERSION)
endif

################################################################################
# Go build details                                                             #
################################################################################

BASE_PACKAGE_NAME := github.com/krancour/go-parrot

LDFLAGS = -w -X $(BASE_PACKAGE_NAME)/version.commit=$(GIT_VERSION) \
	-X $(BASE_PACKAGE_NAME)/version.version=$(EXAMPLES_VERSION)

################################################################################
# Containerized development environment                                        #
################################################################################

DEV_IMAGE := quay.io/deis/lightweight-docker-go:v0.3.0

DOCKER_CMD := docker run \
	--rm \
	-e CGO_ENABLED=0 \
	-e SKIP_DOCKER=true \
	-v $$(pwd):/go/src/$(BASE_PACKAGE_NAME) \
	-v $$(pwd)/.modcache:/go/pkg/mod \
	-w /go/src/$(BASE_PACKAGE_NAME) -it $(DEV_IMAGE)

################################################################################
# Docker images we build and publish                                           #
################################################################################

BASE_IMAGE_NAME        = go-parrot

RC_IMAGE_NAME          = $(DOCKER_REPO)$(BASE_IMAGE_NAME):$(GIT_VERSION)
RC_MUTABLE_IMAGE_NAME  = $(DOCKER_REPO)$(BASE_IMAGE_NAME):canary

REL_IMAGE_NAME         = $(DOCKER_REPO)$(BASE_IMAGE_NAME):$(REL_VERSION)
REL_MUTABLE_IMAGE_NAME = $(DOCKER_REPO)$(BASE_IMAGE_NAME):latest

################################################################################
# Utility targets                                                              #
################################################################################

# Allow developers to step into the containerized development environment--
# unconditionally requires docker
.PHONY: dev
dev:
	$(DOCKER_CMD) bash

# Install/update dependencies
.PHONY: dep
dep:
ifdef SKIP_DOCKER
	scripts/dep.sh
else
	$(DOCKER_CMD) scripts/dep.sh
endif

.PHONY: verify
verify:
ifdef SKIP_DOCKER
	scripts/verify.sh
else
	$(DOCKER_CMD) scripts/verify.sh
endif

.PHONY: linux-stream
linux-stream:
	vlc live-stream.sdp

.PHONY: mac-stream
mac-stream:
	/Applications/VLC.app/Contents/MacOS/VLC live-stream.sdp

.PHONY: land
land: build
	docker run --rm $(RC_IMAGE_NAME) emergency-landing

################################################################################
# Tests                                                                        #
################################################################################

# Executes unit tests
.PHONY: test
test:
ifdef SKIP_DOCKER
	$(UNIT_TEST_CMD) scripts/test.sh
else
	$(DOCKER_CMD) scripts/test.sh
endif

# Executes an extensive series of lint checks against broker code
.PHONY: lint
lint:
ifdef SKIP_DOCKER
	$(LINT_CMD) scripts/lint.sh
else
	$(DOCKER_CMD) scripts/lint.sh
endif

################################################################################
# Misc                                                                         #
################################################################################

EXAMPLE ?= connect

# Build the binaries and docker image from code, then run the specified binary--
# unconditionally requires docker
.PHONY: run
run: build
	docker run --rm $(RC_IMAGE_NAME) $(EXAMPLE)

################################################################################
# Build / Publish                                                              #
################################################################################

# Build the broker binary and docker image
.PHONY: build
build:
	docker build \
		--build-arg BASE_PACKAGE_NAME='$(BASE_PACKAGE_NAME)' \
		--build-arg LDFLAGS='$(LDFLAGS)' \
		-t $(RC_IMAGE_NAME) \
		.
	docker tag $(RC_IMAGE_NAME) $(RC_MUTABLE_IMAGE_NAME)

# Push release candidate image
.PHONY: push-rc
push-rc: build
	docker push $(RC_IMAGE_NAME)
	docker push $(RC_MUTABLE_IMAGE_NAME)

# Rebuild and push officially released, semantically versioned image with
# semantically versioned binary
.PHONY: push-release
push-release:
ifndef REL_VERSION
	$(error REL_VERSION is undefined)
endif
	@# This pull is a verification that this commit has successfully cleared the
	@# master pipeline.
	docker pull $(RC_IMAGE_NAME)
	docker build \
		--build-arg BASE_PACKAGE_NAME='$(BASE_PACKAGE_NAME)' \
		--build-arg LDFLAGS='$(LDFLAGS)' \
		-t $(REL_IMAGE_NAME) \
		.
	docker tag $(REL_IMAGE_NAME) $(REL_MUTABLE_IMAGE_NAME)
	docker push $(REL_IMAGE_NAME)
	docker push $(REL_MUTABLE_IMAGE_NAME)
