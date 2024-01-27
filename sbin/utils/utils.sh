#!/bin/bash

# Converts a path to be relative to another directory.
relative_path() {
  echo $(python3 -c "import os.path; print(os.path.relpath('$1', '$2'))")
}

# Requests a user to confirm the given prompt ($1).
confirm() {
  read -r -p "$1 " response
  if ! [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    exit 1
  fi
}

# Requests a user to confirm that overwriting
# file ($1) is okay, if it exists.
# Doesn't ask if -y is set.
confirm_overwrite() {
  if test -f $1; then
    if [[ "$2" = "false" ]]; then
      read -r -p "Overwrite $1 with a new file? [y/N] " response
      if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        echo "Removing $1"
        rm $1
      else
        exit 1
      fi
    fi
  fi
}

# Requires that a dotenv named $1 exists at a path ($2).
require_dotenv() {
  if ! test -f "$2"; then
    echo "Expected $1 dotenv at $2 (not found)."
    exit 1
  fi
  echo "Using $1 dotenv: $2"
  . $2
}

# Requires that all env variables named in $@ are set.
require_env() {
  for var in "$@"; do
    if [ -z ${!var+x} ]; then
      echo "$var is required but not set"
      exit 1
    fi
  done
}
