#!/bin/bash -eu

if [ ! -d /Volumes/Contents/a7iii/ ]; then
  echo "Drive is not mounted."
  exit 1
fi

cp -rf /Volumes/1-1/DCIM/* ~/Pictures/tmp/
cd ~/Pictures/tmp
rename-sony-photos-directories
cp -rf ~/Pictures/tmp/* /Volumes/Contents/a7iii/
rm -rf /Volumes/1-1/DCIM/*
rm -rf ~/Pictures/tmp/*
diskutil eject 1-1
