#!/bin/bash
if [[ $# -eq 0 ]]; then
  echo "no test name provided"
  exit 1
fi

# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

###### Process handling ######
trap ctrl_c INT

function cleanup() {
  echo "Cleaning up processes..."
  pgrep geth | xargs kill
  pgrep sidecar | xargs kill
  pgrep magi | xargs kill
  L1_PORT=4012
  if [ -n "$L1_PORT" ]; then
    L1_PORT_PID=$(lsof -i tcp:${L1_PORT} | awk 'NR!=1 {print $2}')
    if [ -n "$L1_WS_PID" ]; then
      echo "Killing proc on $L1_PORT"
      kill $L1_PORT_PID
    fi
  fi
  # Clean up
  $SBIN/clean.sh
}

function ctrl_c() {
  cleanup
  exit
}

##############################

WORKSPACE_DIR=$HOME/.spc/workspaces/active_workspace

# Copy config files to workspace.
spc workspace download --config-path "config/spc?ref=siosw/spc-integration" --name e2e
spc workspace set e2e

PATHS_ENV=$WORKSPACE_DIR/.paths.env
SIDECAR_ENV=$WORKSPACE_DIR/.sidecar.env
reqdotenv "paths" $PATHS_ENV
reqdotenv "sidecar" $SIDECAR_ENV

# Start L1
yes | $SBIN/generate_secrets.sh -dj
echo "starting l1"
yes | $SBIN/start_l1.sh -d &>$WORKSPACE_DIR/l1.out &

echo "waiting for host"

# Parse url into host:port
L1_HOST_AND_PORT=${L1_ENDPOINT#*://}
echo $L1_HOST_AND_PORT
# Wait for services
$SBIN/wait-for-it.sh -t 60 $L1_HOST_AND_PORT | sed "s/^/[WAIT] /"
echo "L1 endpoint is available"

DEPLOYMENTS_ENV=$WORKSPACE_DIR/.deployments.env
until [ -f $DEPLOYMENTS_ENV ]; do
  echo "waiting for L1 to be fully deployed..."
  sleep 4
done
reqdotenv "deployments" $DEPLOYMENTS_ENV

# Start sp-geth
$SBIN/start_sp_geth.sh -c &>$WORKSPACE_DIR/sp_geth.out &
sleep 2

# Start sp-magi
$SBIN/start_sp_magi.sh &>$WORKSPACE_DIR/sp_magi.out &
sleep 2

# Start sidecar
$SBIN/start_sidecar.sh &>$WORKSPACE_DIR/sidecar.out &
sleep 2

cd $CONTRACTS_DIR
echo "Running test: $1"
echo $(pwd)
# Run testing script
case $1 in
transactions)
  npx hardhat run scripts/e2e/test_transactions.ts
  RESULT=$?
  ;;
deposit)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_deposit_eth.ts
  RESULT=$?
  ;;
withdraw)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_withdraw_eth.ts
  RESULT=$?
  ;;
erc20)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_erc20.ts
  RESULT=$?
  ;;
*)
  echo "unknown test"
  RESULT=1
  ;;
esac

cleanup
exit $RESULT
