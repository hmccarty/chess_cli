#!/bin/bash

DEV_DIR="$(readlink -f $(dirname $0)/../)"

docker run -it --rm \
        -e USER \
        -v $DEV_DIR:/app \
        --name gochess-dev \
        gochess \
        go mod tidy