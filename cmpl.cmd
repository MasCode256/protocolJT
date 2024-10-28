@echo off

go build -o methods\echo\echo.exe methods\echo\echo.go
go build -o code\client\client.exe code\client\client.go
go build -o code\server\server.exe code\server\server.go
go build -o code\tracker\tracker.exe code\tracker\tracker.go

taskkill /IM server.exe /F
taskkill /IM tracker.exe /F

start /b methods\echo\echo.exe
start /b code\server\server.exe
start /b code\tracker\tracker.exe