
@echo off

if not "%GO_INIT%"=="1" (
set GO_INIT=1
set GOPATH=%GOPATH%;%cd%
)

set GOARCH=386
set CGO_ENABLED=1

go build -ldflags "-s -w" -o bin/release/lanDDNS.exe
