CC=go
RM=rm
MV=mv


SOURCEDIR=./yeti
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
APP=Ympact

VERSIONCOMMENT:=$(shell egrep -m1 "Appversion" yeti/ui/extract_register_ui.go | cut -d= -f2 | tr -d "\"" | tr " " "_")
suffix=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2 | sed 's/[0-9.]//g')
snapshot=$(shell date +%FT%T)
VERSION="1.0-"$(VERSIONCOMMENT)
DATE:=$(shell date +"%Y-%m-%d")

ifeq ($(suffix),rc)
	appversion=$(VERSION)$(snapshot)
else 
	appversion=$(DATE)
endif 

.DEFAULT_GOAL:=build


build:
	@echo "Get the dependencies"
	go mod tidy 
	(make build-darwin)
	(make build-windows)


build-darwin:
	@echo "Compilation for macos"
	fyne package -os darwin -icon ./Icon.png -name $(APP) -source-dir $(SOURCEDIR)/
	zip -r $(APP)-$(appversion)-macos.zip $(APP).app

build-windows:
	@echo "Compilation for windows"
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -o $(APP).exe $(SOURCEDIR)/
	zip $(APP)-$(appversion)-windows.zip $(APP).exe ./dll/opengl32.dll

clean:
	@echo "Cleaning all *.zip archives."
	rm -f $(APP)*.zip
	@echo "Cleaning all binaries."
	rm -fr $(APP)-* $(APP).exe

deps: get-linter get-vulncheck
	@echo "Getting tools..."

get-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

get-vulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest

lint:
	@echo "Lint the whole project"
	golangci-lint run --timeout 5m ./...

vulncheck:
	govulncheck ./...
