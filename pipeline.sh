#!/bin/bash
changes=$(git diff --name-only | cut -d'/' -f2| uniq)
services=()
for dir in changes;do
  if [ -f $dir/Dockerfile ];then
    services+=("$dir")
  fi
done
echo "services:$services" 
services >> $GITHUB_OUTPUT
