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

################################################################################
# Containerized development environment                                        #
################################################################################

DEV_IMAGE := quay.io/deis/lightweight-docker-go:v0.3.0

DOCKER_CMD := docker run \
	-it \
	--rm \
	-e CGO_ENABLED=0 \
	-e SKIP_DOCKER=true \
	-v $$(pwd):/go/src/$(BASE_PACKAGE_NAME) \
	-v $$(pwd)/.modcache:/go/pkg/mod \
	-w /go/src/$(BASE_PACKAGE_NAME) $(DEV_IMAGE)

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
