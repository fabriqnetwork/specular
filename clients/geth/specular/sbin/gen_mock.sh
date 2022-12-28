#!/usr/bin/env bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
SPECULAR_DIR=$SBIN/../


SRC=(
    $SPECULAR_DIR/rollup/services/backend.go
    $SPECULAR_DIR/proof/api.go
)

for path in ${SRC[@]}; do 
    fname=`basename $path`
    dir=`dirname $path`
    dir=`cd $dir; pwd`

    src=$dir/$fname
    dest=$dir/mock/"${fname%%.*}"_mock.go
    echo "Generating mock for $src at $dest"
    mockgen -package mock -source=$src > $dest
done
