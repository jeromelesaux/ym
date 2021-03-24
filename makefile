CC=go
RM=rm
MV=mv


SOURCEDIR=./extract
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION:="1.0-hexa-emptyselection"
suffix=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2 | sed 's/[0-9.]//g')
snapshot=$(shell date +%FT%T)

ifeq ($(suffix),rc)
	appversion=$(VERSION)$(snapshot)
else 
	appversion=$(VERSION)
endif 

.DEFAULT_GOAL:=build


build: 
	@echo "Get the dependencies" 
	go install fyne.io/fyne/v2/cmd/fyne
	@echo "Compilation for linux"
	GOOS=linux && CGO_ENABLED=1 && go build ${LDFLAGS} -o YeTi $(SOURCEDIR)/
	zip YeTi-$(appversion)-linux.zip YeTi
	@echo "Compilation for macos"
	fyne package -os darwin -icon  $(SOURCEDIR)/icon/YeTi.png -name YeTi -sourceDir $(SOURCEDIR)/
	zip -r YeTi-$(appversion)-macos.zip YeTi.app
	@echo "Compilation for windows"
	export GOOS=windows  && export GOARCH=386  && export CGO_ENABLED=1 && export CC=i686-w64-mingw32-gcc && go build ${LDFLAGS} -o YeTi.exe $(SOURCEDIR)/
	zip YeTi-$(appversion)-windows.zip YeTi.exe

clean:
	@echo "Cleaning all *.zip archives."
	rm -f YeTi*.zip
	@echo "Cleaning all binaries."
	rm -fr YeTi*