@ECHO OFF

del phone.exe
go build -mod=vendor -ldflags="-s -w" -o phone.exe phone.go

PAUSE
