#!/usr/bin/env bash

# AVOID INVOKING THIS SCRIPT DIRECTLY -- USE `make verify`

# This script should only be used by CI. It provides a sanity check that
# dependency resolution didn't find any NEW dependencies that haven't been
# recorded in go.mod and go.sum.

set -euo pipefail

scripts/dep.sh

if [ -z "$(git status --porcelain)" ]; then
  exit 0
else
  cat << EOF

After dependency resolution, the working directory contains uncommitted changes!

This probably indicates that dependencies NOT recorded in go.mod and go.sum have
been discovered. Without these recorded, future repeatability of the build
cannot be guaranteed.

To resolve this, please run `make dep` then commit the updated go.mod, go.sum,
and vendor/ directory.

EOF
  exit 1
fi
