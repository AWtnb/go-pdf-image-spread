@echo off
cd /d %~dp0

echo;
call converter.exe --vertical --singleTop

echo;
echo;
echo finished!
echo;
pause