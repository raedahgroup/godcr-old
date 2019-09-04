module github.com/raedahgroup/godcr

go 1.12

require github.com/raedahgroup/godcr/cmd v0.0.0-00010101000000-000000000000 // indirect

replace (
	github.com/raedahgroup/godcr/app => ./app
	github.com/raedahgroup/godcr/cli => ./cli
	github.com/raedahgroup/godcr/cmd => ./cmd
	github.com/raedahgroup/godcr/fyne => ./fyne
	github.com/raedahgroup/godcr/nuklear => ./nuklear
	github.com/raedahgroup/godcr/terminal => ./terminal
	github.com/raedahgroup/godcr/web => ./web
)
