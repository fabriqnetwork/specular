#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
IS_SEQUENCER=1 $SBIN/start_node.sh "$@"
