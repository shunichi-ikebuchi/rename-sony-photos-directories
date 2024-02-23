# rename-sony-photos-directores

This project contains a program and shell script to change [date format folder](https://www.sony.jp/ServiceArea/impdf/pdf/44879440M.w-JP/jp/contents/TP0000220296.html) of Sony digital cameras to yyyy-mm-dd format. dd format.

# Prerequisites

- macOS
  - The shell scripts included in this repository contain commands valid only on macOS
- The name of the SD card where you want to store the photos must be `1-1`.
- The path of the SD card where you want to store the photos must be `/Volumes/1-1`.
- The name of the SD card for backup must be `1-2
- The path of the SD card you want to back up must be `/Volumes/1-2` The name of the SD card you want to back up must be `1-2



## How to use

## Preparation

Go to `go install` and place `rename-sony-photos-directores` in a location with PATH.



## Rename the folder and copy the photos.

``` .
. /copy-rename-and-delete.sh
````

## delete photos on SD card for backup ```` .
```
. /delete.sh
```