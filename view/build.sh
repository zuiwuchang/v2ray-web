#!/bin/bash
DirRoot=`cd $(dirname $BASH_SOURCE) && pwd`
cd "$DirRoot" && ng build --configuration production --base-href /view/ --localize
