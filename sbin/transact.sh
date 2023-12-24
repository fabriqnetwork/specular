#!/bin/bash
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
ROOT_DIR=$SBIN/..

SP_GETH_ENV=".sp_geth.env"
if ! test -f "$SP_GETH_ENV"; then
  echo "Expected dotenv at $SP_GETH_ENV (does not exist)."
  exit
fi
echo "Using dotenv: $SP_GETH_ENV"
. $SP_GETH_ENV

# VALUE=0.0001
# for i in $(seq 1 $END);
#     cast send --async \
#         --rpc-url http://$ADDRESS:$HTTP_PORT \
#         --chain $NETWORK_ID \
#         --private-key `cat $1` \
#         --value "$VALUE"ether \
#         0x0000000000000000000000000000000000000000
# done
