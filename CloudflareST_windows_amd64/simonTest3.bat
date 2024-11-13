@echo off
:: 获取当前时间戳
setlocal enabledelayedexpansion
set "timestamp=%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%"
set "timestamp=%timestamp: =0%"

:: 设置文件名
set "filename=result3\result%timestamp%.csv"

:: 打印生成的文件名
:: echo 生成的文件名: %filename%
:: -dd 禁止下载测速
:: 运行 CloudflareST.exe
main.exe -f IPlist3.txt -url https://speed.cloudflare.com/__down?bytes=52428800 -o %filename%  -p 20 -dn 1000 -n 10 -tp 443 -httping -httping-code 404
::main.exe -f IPlist3.txt -url https://www.speedtest.net -o %filename%  -p 20 -dn 1000 -dd -httping -n 10 -tp 443
:: CloudflareST.exe -f IPlist3.txt -url https://www.speedtest.net -o %filename%  -p 20 -dn 1000 -dd 

:: https://speed.cloudflare.com/__down?bytes=104857600