@echo off

if not "%GO_INIT%"=="1" (
set GO_INIT=1
set GOPATH=%GOPATH%;%cd%
)

set GOARCH=386
set CGO_ENABLED=1

go run "src/test/ComTest.go"
