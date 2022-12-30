@echo off
cd /d %~dp0

echo;
call converter.exe --singleTop

echo;
echo;
echo finished!
echo;
pause