@echo off
cd /d %~dp0

echo;
call converter.exe --vertical --singleTop

echo;
echo;

if %errorlevel% equ 0 echo finished!
if %errorlevel% equ 1 echo ERROR!
echo;
pause