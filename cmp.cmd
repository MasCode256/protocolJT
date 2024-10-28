@echo off

go build -o code\client\client.exe code\client\client.go
go build -o code\server\server.exe code\server\server.go
go build -o code\tracker\tracker.exe code\tracker\tracker.go