#!/bin/bash -ex

for d in $(find * -name 'go.mod' -type f | sed -r 's|/[^/]+$||' | sort -u); do
    go work use $d
done