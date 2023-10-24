#!/bin/bash

version_file=VERSION

# Check if the version file exists
if [[ ! -f $version_file ]]; then
    echo "$version_file does not exist"
    exit 1
fi

# Get the version from the file
version=$(cat $version_file)

# Tag the commit with the version number
git tag "v$version"
