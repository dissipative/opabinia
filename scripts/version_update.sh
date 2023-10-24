#!/bin/bash

# Path to the version file
version_file=VERSION

# Check if the version file exists
if [[ ! -f $version_file ]]; then
    echo "$version_file does not exist"
    exit 1
fi

# Check if the VERSION file is staged for commit
if git diff --cached --name-only --quiet -- $version_file
then
    # If the VERSION file is not staged for commit, increment the version
    # Read the version from the file
    version=$(cat $version_file)

    # Split the version by dots
    IFS='.' read -ra version_parts <<< "$version"

    # Increment the minor version
    ((version_parts[2]++))

    # Reassemble the version
    new_version="${version_parts[0]}.${version_parts[1]}.${version_parts[2]}"

    # Write the new version to the file
    echo "$new_version" > $version_file

    # Add the updated version file to the commit
    git add $version_file
fi

