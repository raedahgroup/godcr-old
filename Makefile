.POSIX:

.PHONY: build install

# Target rules to build all binaries
build: godcr-fyne

# Target rules to build and install all binaries
install: install-fyne

godcr-fyne: fyne/fyne-packr.go fyne/packrd/packed-packr.go
	go build ./cmd/godcr-fyne

install-fyne: fyne/fyne-packr.go fyne/packrd/packed-packr.go
	go install ./cmd/godcr-fyne

fyne/fyne-packr.go fyne/packrd/packed-packr.go: fyne/display.go fyne/icons.go
	cd fyne && packr2
