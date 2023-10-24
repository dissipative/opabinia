#!/bin/bash

echo "Fixing code files..."
make hooks.style
echo "Updating version if necessary..."
make hooks.version-update
