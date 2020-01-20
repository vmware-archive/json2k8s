#!/bin/bash

set -e -x -u

./hack/build.sh

cat examples/file1.json | ./json2k8s
cat examples/file1.json | ./json2k8s -
./json2k8s examples/file1.json
./json2k8s examples/file2.json
./json2k8s examples/*.json

echo ALL SUCCESS
