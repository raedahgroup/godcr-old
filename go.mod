module github.com/raedahgroup/godcr

go 1.12

require github.com/raedahgroup/godcr/cmd/godcr-fyne v0.0.0-00010101000000-000000000000 // indirect

replace (
	github.com/raedahgroup/godcr/cmd/godcr-fyne => ./cmd/godcr-fyne
	github.com/raedahgroup/godcr/fyne => ./fyne
)
