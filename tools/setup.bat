@echo off
REM
REM Usage:  C:\> setup.bat expand_dir
REM

for /F "tokens=1 delims=;" %%a in ('go env GOPATH') do set go_path=%%a

xcopy /e %1\Resources %go_path%\bin
xcopy /e %1\Release   %go_path%\bin

set CGO_CFLAGS=-I%1
set CGO_LDFLAGS=-L%go_path%\bin -lcef
