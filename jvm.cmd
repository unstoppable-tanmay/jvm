@echo off
REM jvm Command Prompt (cmd.exe) wrapper
REM Sourced via DOSKEY or placed in PATH as jvm.cmd

IF "%1"=="use" (
    "%~dp0jvm.exe" %*
    IF ERRORLEVEL 0 (
        FOR /F "tokens=*" %%i IN ('"%~dp0jvm.exe" env --shell cmd') DO %%i
    )
) ELSE (
    "%~dp0jvm.exe" %*
)
