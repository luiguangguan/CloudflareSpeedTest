@echo off
:: 获取当前时间戳
setlocal enabledelayedexpansion
set "timestamp=%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%"
set "timestamp=%timestamp: =0%"

:: 设置文件名
set "filename=result2\result%timestamp%.txt"

:: 打印生成的文件名
:: echo 生成的文件名: %filename%

:: 运行 CloudflareST.exe
CloudflareST.exe -f IPlist2.txt -url https://speed.cloudflare.com/__down?bytes=52428800 -n 10 -o %filename%  -p 20 -dn 1000

:: https://speed.cloudflare.com/__down?bytes=104857600