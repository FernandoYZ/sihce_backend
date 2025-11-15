@echo off
REM --- Cargar secretos ---
call env.bat

REM --- Ejecutar la aplicación Go ---
go run .

pause
