# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

reqdotenv ""

# generate_wallet.sh

# modify: genesis_config.json accs
# modify: base_rollup.json
