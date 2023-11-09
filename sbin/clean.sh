#!/bin/bash
SBIN=`dirname $0`
$SBIN/clean_sp_geth.sh
$SBIN/clean_deployment.sh

echo "Removing dotenv files..."
rm -f .contracts.env
rm -f .genesis.env
rm -f .sp_geth.env
rm -f .sp_magi.env
rm -f .sidecar.env
echo "Done."
