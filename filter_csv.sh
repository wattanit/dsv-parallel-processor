#!/bin/zsh

SPEC_DIR=$(dirname $(realpath $1 ))
DATA_DIR=$(realpath $2 )
echo "spec file directory = $SPEC_DIR"
echo "input file directory = $DATA_DIR"

docker run -it --rm -v $SPEC_DIR:/spec -v $DATA_DIR:/data dsv-parallel-processor:0.1 /spec/$1