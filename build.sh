#!/bin/bash
#Program:
#       golang build scripts
#
#Email:
#       zuiwuchang@gmail.com
DirRoot=`cd $(dirname $BASH_SOURCE) && pwd`
Target=v2ray-web
TestItems=(
)
function check(){
	if [ "$1" != 0 ] ;then
		exit $1
	fi
}
function mkDir(){
	mkdir -p "$1"
	check $?
}

function newFile(){
	echo "$2" > "$1"
	check $?
}
function writeFile(){
	echo "$2" >> "$1"
	check $?
}
function createGoVersion(){
	mkDir "$DirRoot/version"
	filename="$DirRoot/version/version.go"
	package=version


	tag=`git describe`
	if [ "$tag" == '' ];then
		tag="[unknown tag]"
	fi

	commit=`git rev-parse HEAD`
	if [ "$commit" == '' ];then
		commit="[unknow commit]"
	fi
	
	date=`date +'%Y-%m-%d %H:%M:%S'`

	echo $tag $commit
	echo $date

	newFile $filename	"package $package"
	writeFile $filename	''
	writeFile $filename	'// Tag git tag'
	writeFile $filename	"const Tag = \`$tag\`"
	writeFile $filename	'// Commit git commit'
	writeFile $filename	"const Commit = \`$commit\`"
	writeFile $filename	'// Date build datetime'
	writeFile $filename	"const Date = \`$date\`"
}
function buildSource(){
	echo statik -src="$DirRoot/view/dist/view/$1" -dest "$DirRoot/assets/$1"  -ns "$1" -f
	rm "$DirRoot/assets/$1" -rf
	mkdir "$DirRoot/assets/$1" -p
	statik -src="$DirRoot/view/dist/view/$1" -dest "$DirRoot/assets/$1"  -ns "$1" -f
}

function DisplayHelp(){
	echo "help                       : display help"
	echo "l/linux   [r/d] [t/tar]    : build for linux"
	echo "d/darwin  [r/d] [t/tar]    : build for darwin"
	echo "w/windows [r/d] [t/tar]    : build for windows"
	echo "s/source                   : build assets"
	echo "t/test                     : run go test"
}
case $1 in
	l|linux)
		export GOOS=linux
		export CGO_ENABLED=0

		createGoVersion
		if [[ $2 == d ]]; then
			Target="$Target"d
			echo go build -tags=jsoniter -o "$DirRoot/bin/$Target"
			cd "$DirRoot" && go build -tags=jsoniter -o "$DirRoot/bin/$Target"
		else
			echo go build -tags=jsoniter -ldflags "-s -w" -o "$DirRoot/bin/$Target"
			cd "$DirRoot" && go build -tags=jsoniter -ldflags "-s -w" -o "$DirRoot/bin/$Target"
		fi
		check $?

		if [[ $3 == tar || $3 == t ]]; then
			dst=linux.amd64.tar.xz
			if [[ $GOARCH == 386 ]];then
				dst=linux.386.tar.xz
			fi
			cd "$DirRoot/bin" && tar -Jcvf $dst "$Target" "$Target.jsonnet" \
				geoip.dat geosite.dat \
				v2ray-web.service \
				run.sh
		fi
	;;

	d|darwin)
		export GOOS=darwin
		export CGO_ENABLED=0

		createGoVersion
		if [[ $2 == d ]]; then
			Target="$Target"d
			echo go build -tags=jsoniter -o "$DirRoot/bin/$Target"
			cd "$DirRoot" && go build -tags=jsoniter -o "$DirRoot/bin/$Target"
		else
			echo go build -tags=jsoniter -ldflags "-s -w" -o "$DirRoot/bin/$Target"
			cd "$DirRoot" && go build -tags=jsoniter -ldflags "-s -w" -o "$DirRoot/bin/$Target"
		fi
		check $?

		if [[ $3 == tar || $3 == t ]]; then
			dst=darwin.amd64.tar.xz
			if [[ $GOARCH == 386 ]];then
				dst=darwin.386.tar.xz
			fi
			cd "$DirRoot/bin" && tar -Jcvf $dst "$Target" "$Target.jsonnet" \
				geoip.dat geosite.dat \
				v2ray-web.service \
				run.sh
		fi
	;;

	w|windows)
		export GOOS=windows
		export CGO_ENABLED=0

		createGoVersion
		if [[ $2 == d ]]; then
			Target="$Target"d
			echo go build -tags=jsoniter -o "$DirRoot/bin/$Target.exe"
			cd "$DirRoot" && go build -tags=jsoniter -o "$DirRoot/bin/$Target.exe"
		else
			echo go build -tags=jsoniter -ldflags "-s -w" -o "$DirRoot/bin/$Target.exe"
			cd "$DirRoot" && go build -tags=jsoniter -ldflags "-s -w" -o "$DirRoot/bin/$Target.exe"
		fi
		check $?

		if [[ $3 == tar || $3 == t ]]; then
			dst=windows.amd64.tar.xz
			if [[ $GOARCH == 386 ]];then
				dst=windows.386.tar.xz
			fi
			cd "$DirRoot/bin" && tar -Jcvf $dst "$Target.exe" "$Target.jsonnet" \
				geoip.dat geosite.dat \
				v2ray-web-service.xml v2ray-web-service.exe \
				run.bat install.bat
		fi
	;;

	s|source)
		cp "$DirRoot/view/dist/view/en-US/3rdpartylicenses.txt" "$DirRoot/static/3rdpartylicenses.txt"

		echo statik -src="$DirRoot/static" -dest "$DirRoot/assets/static"  -ns static -f
		rm "$DirRoot/assets/static" -rf
		mkdir "$DirRoot/assets/static" -p
		statik -src="$DirRoot/static" -dest "$DirRoot/assets/static"  -ns static -f

		buildSource zh-Hant
		buildSource zh-Hans
	;;

	t|test)
		for i in ${!TestItems[@]}
		do
			cd "$DirRoot/${TestItems[i]}" && go test
		done
	;;

	*)
		DisplayHelp
	;;
esac