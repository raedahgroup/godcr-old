module github.com/raedahgroup/godcr

go 1.12

require (
	gioui.org v0.0.0-20191022064259-288460452184 // indirect
	gioui.org/ui v0.0.0-20190926171558-ce74bc0cbaea // indirect
	github.com/aarzilli/nucular v0.0.0-20191004125635-0f0b2bda58e2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/raedahgroup/dcrseedgen v0.0.0-20190903182641-0cb9eeab9be9 // indirect
	github.com/raedahgroup/godcr/cmd/godcr-fyne v0.0.0-00010101000000-000000000000 // indirect
	github.com/raedahgroup/godcr/cmd/godcr-terminal v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	github.com/raedahgroup/godcr/cmd/godcr-fyne => ./cmd/godcr-fyne
	github.com/raedahgroup/godcr/cmd/godcr-terminal => ./cmd/godcr-terminal
	github.com/raedahgroup/godcr/fyne => ./fyne
	github.com/raedahgroup/godcr/terminal => ./terminal
)

require github.com/raedahgroup/godcr/cmd/godcr-fyne v0.0.0-00010101000000-000000000000 // indirect

replace (
	github.com/raedahgroup/godcr/cmd/godcr-fyne => ./cmd/godcr-fyne
	github.com/raedahgroup/godcr/cmd/godcr-terminal => ./cmd/godcr-terminal
	github.com/raedahgroup/godcr/fyne => ./fyne
	github.com/raedahgroup/godcr/terminal => ./terminal
)
