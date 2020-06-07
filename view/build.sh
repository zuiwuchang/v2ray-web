#!/bin/bash
DirRoot=`cd $(dirname $BASH_SOURCE) && pwd`
cd "$DirRoot" && ng build --prod  --base-href "/view/" --localize --lazy-modules
