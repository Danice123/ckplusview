#!/bin/sh
GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc fyne package -os windows -icon icon.png