#!/bin/bash

set -e -x -u

go fmt ./cmd/...

# makes builds reproducible
export CGO_ENABLED=0
repro_flags="-ldflags=-buildid= -trimpath"

go build $repro_flags -o json2k8s ./cmd/...
