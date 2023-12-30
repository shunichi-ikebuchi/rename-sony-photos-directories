# This script deletes all files from the SD card and ejects it.
#!/bin/bash -eu

if [ ! -d /Volumes/1-2/ ]; then
  echo "Drive is not mounted."
  exit 1
fi

rm -rf /Volumes/1-2/DCIM/*
diskutil eject 1-2
