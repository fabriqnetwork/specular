#!/bin/bash
cd /specular/workspace
shopt  -s dotglob
echo *

# Run the main container command.
exec "$@"
