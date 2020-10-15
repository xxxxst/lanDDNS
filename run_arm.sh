#!/bin/bash

rootDir=$(pwd)
if [ -z $GO_INIT ];then
	export GO_INIT=1
	export GOPATH=$GOPATH:$rootDir
fi

export GOARCH=arm
export CGO_ENABLED=1

go run -tags=cldebug main.go

