.PHONY: all dockerize-db clean

ifeq ($(OS),Windows_NT)
    MAKE = "C:/MinGW/msys/1.0/bin/make.exe"
endif

all: 
	"${MAKE}" -C cmd
	"${MAKE}" -C db

dockerize-db:
	"${MAKE}" -C db

clean:
	"${MAKE}" -C cmd clean