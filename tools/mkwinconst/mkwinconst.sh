#!/bin/bash
cat /usr/include/w32api/winuser.h | grep -e '^#define[ \t]*WS_' |\
 sed -e 's/^#define[ \t]*\(WS_[_A-Z]*\)[ \t].*$/\1/' |\
go run mkwinconst.go
