#!/bin/bash

for directory in `find . -type f -name '*.go' -exec dirname {} + | sort -u`; do
    go fmt $directory
done

