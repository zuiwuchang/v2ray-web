#!/bin/bash
DirRoot=`cd $(dirname $BASH_SOURCE) && pwd`
opencc -i "$DirRoot/src/locale/zh-Hant.xlf" -o "$DirRoot/src/locale/zh-Hans.xlf" -c t2s.json
