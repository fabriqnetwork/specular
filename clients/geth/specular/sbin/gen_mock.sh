#!/usr/bin/env bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
SPECULAR_DIR=$SBIN/../
GETH_DIR=$SBIN/../../go-ethereum

echo "This is SBIN: $SBIN"
echo "This is GETH_DIR: $GETH_DIR"

SRC=(
    # $SPECULAR_DIR/rollup/services/backend.go
    # $SPECULAR_DIR/proof/api.go
    #$GETH_DIR/ethclient/ethclient.go
    $SPECULAR_DIR/rollup/services/testutils/client.go
)

for path in ${SRC[@]}; do 
    fname=`basename $path`
    echo "This is fname: $fname"
    dir=`dirname $path`
    dir=`cd $dir; pwd`

    echo "This is dir: $dir"

    src=$dir/$fname
    echo "This is src: $src"
    #src=$
    # dest=$dir/mock/"${fname%%.*}"_mock.go
    dest=$SPECULAR_DIR/rollup/services/mock/"${fname%%.*}"_mock.go
    echo "Generating mock for $src at $dest"
    mockgen -package mock -source=$src > $dest
done
