#!/bin/bash

# Source and destination directories
src_dir="../contracts/typechain-types/"
dest_dir="./src/typechain-types/"

# Check if source directory exists
if [ ! -d "$src_dir" ]; then
  echo "Error: Source directory '$src_dir' does not exist."
  exit 1
fi

# Create destination directory if it doesn't exist
if [ ! -d "$dest_dir" ]; then
  mkdir -p "$dest_dir"
fi

# Copy the contents from source to destination
cp -r "$src_dir"* "$dest_dir"

echo "Successfully copied files from '$src_dir' to '$dest_dir'"
