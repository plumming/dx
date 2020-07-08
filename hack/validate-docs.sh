#!/bin/bash

set -e
set -u

./build/linux/chilly docs
git status

CHANGED=$(git ls-files --modified --others --exclude-standard)
if [ "${CHANGED}" == "" ];
then
  echo "All generated docs up to date";
  git diff
else
  echo "Doc generation is out of date";
  echo "${CHANGED}"
  git diff
  exit 1
fi
