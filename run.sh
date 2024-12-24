#!/bin/bash -ex

/usr/bin/clear

# /usr/bin/go \
#     run ./main.go \
#     --help

/usr/bin/go \
    run ./main.go \
    --arch arm64 \
    --version 1.96.0
