#!/bin/bash
SBIN=$(dirname $0)
WORKSPACE_DIR=$HOME/.spc/workspaces/active_workspace

$SBIN/clean_sp_geth.sh
$SBIN/clean_deployment.sh

echo "Removing dotenv files..."
rm -f $WORKSPACE_DIR/.contracts.env
rm -f $WORKSPACE_DIR/.genesis.env
rm -f $WORKSPACE_DIR/.sp_geth.env
rm -f $WORKSPACE_DIR/.sp_magi.env
rm -f $WORKSPACE_DIR/.sidecar.env
rm -f $WORKSPACE_DIR/.paths.env
echo "Done."

echo "Removing $JWT_SECRET_PATH"
rm -f $JWT_SECRET_PATH
