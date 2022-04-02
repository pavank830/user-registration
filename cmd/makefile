.PHONY: all docker clean

all:
	go mod tidy
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=0 go build  -o ../bin/userRegistration.exe
else
	CGO_ENABLED=0 go build  -o ../bin/userRegistration
endif

clean:
ifeq ($(OS),Windows_NT)
	-powershell.exe -ExecutionPolicy Bypass -Command 'Remove-Item -Path ../bin/userRegistration.exe -ErrorAction Continue'
else
	-rm -f ../bin/userRegistration
endif