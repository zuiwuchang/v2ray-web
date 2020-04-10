#!/bin/bash
#Program:
#       vscode 的 gomod 卡到根本無法使用
#       好在 gomod 提供了 vendor 指令 可以自動創建 vendor 檔案夾 並將依賴庫 copy 到此
#       所以 只需要 使用 gomod vendor 創建 vendor 檔案夾，此時項目只要存在與 GOPATH 中 即可在 禁用gomod時 完全正常的工作
#
#       於是乎
#       1. 在命令行中使用 go mod 管理 項目依賴
#       2. 重設 GOPATH 將當前項目加入到 GOPATH 中 ，禁用 GOMOD ，以新環境啓動 vscode
#       
#       這個腳本 正是爲了完成此功能而寫
#       
#History:
#       2019-03-16 king first release
#Email:
#       zuiwuchang@gmail.com


function ShowHelp(){
	echo "h/help     : show help"
	echo "o/open     : open vscode"
	echo "v/vendor   : update vendor"
	echo "n/new      : new project"
}

function Open(){
	dir=$1
	filename=$dir/go.mod
	echo $filename
	while read -r line
	do
		if echo "$line" | grep -q '^module'; then
			module=`echo $line | awk '{ print $2 }'`
			break
		fi
	done < "$filename"
	CheckError $?
	echo module $module
	gopath=$dir/..
	while [[ $module != "." ]]
	do
		module=`dirname $module`
		gopath=$gopath/..
	done
	gopath=`cd $gopath && pwd`
	CheckError $?
	echo $gopath
	echo GOPATH $gopath
	export GOPATH=$gopath:$GOPATH

	export GO111MODULE=off

	code $dir
}
function CheckError(){
	if [[ $1 != 0 ]];then
		exit $1
	fi
}
case $1 in
	o|open)
		if [[ $2 == "" ]]; then
			Open `cd $(dirname $BASH_SOURCE) && pwd`
		else
			Open $2
		fi
	;;
	
	v|vendor)
		if [[ $2 == "" ]]; then
			dir=`cd $(dirname $BASH_SOURCE) && pwd`
			CheckError $?
			export GO111MODULE=on
			`cd $dir && go mod vendor`
		else
			`cd $2 && go mod vendor`
		fi
	;;

	n|new)
		mkdir $2 -p
		CheckError $?
		`cd $2 && go mod init $2`
		CheckError $?
		dir=`cd $(dirname $BASH_SOURCE) && pwd`
		CheckError $?
		cp $dir/$(basename $BASH_SOURCE) $2/
	;;
	
	h|help)
		ShowHelp
	;;

	*)
		ShowHelp
		exit 1
	;;
esac
