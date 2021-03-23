## YeTi or YM file extractor from ImPact

Project preparation : 
- ```go get fyne.io/fyne/v2```
- ```go install fyne.io/fyne/v2/cmd/fyne```

to compile for windows : 
```export  GOOS="windows"  && export GOARCH="386"  && export  CGO_ENABLED="1" && export CC="i686-w64-mingw32-gcc" 
go build -o YeTi.exe
```
to compile and package for macos : 
```
fyne package -os darwin -icon  icon/YeTi.png -name YeTi
```
