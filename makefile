CC=go
RM=rm
MV=mv


SOURCEDIR=./extract
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSIONCOMMENT:=$(shell egrep -m1 "Appversion" extract/ui/extract_register_ui.go | cut -d= -f2 | tr -d "\"" | tr " " "_")
suffix=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2 | sed 's/[0-9.]//g')
snapshot=$(shell date +%FT%T)
VERSION="1.0-"$(VERSIONCOMMENT)

ifeq ($(suffix),rc)
	appversion=$(VERSION)$(snapshot)
else 
	appversion=$(VERSION)
endif 

.DEFAULT_GOAL:=build


build:
	@echo "Get the dependencies"
	go mod tidy 
	@echo "Compilation for macos"
	rm -f $(SOURCEDIR)/extract
	fyne package -os darwin -icon  $(SOURCEDIR)/icon/YeTi.png -name YeTi -sourceDir $(SOURCEDIR)/
	zip -r YeTi-$(appversion)-macos.zip YeTi.app
	@echo "Compilation for windows"
	export GOOS=windows && export GOARCH=386 && export CGO_ENABLED=1 && export CC=i686-w64-mingw32-gcc && go build ${LDFLAGS} -o YeTi.exe $(SOURCEDIR)/
	zip YeTi-$(appversion)-windows.zip YeTi.exe

clean:
	@echo "Cleaning all *.zip archives."
	rm -f YeTi*.zip
	@echo "Cleaning all binaries."
	rm -fr YeTi*

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
