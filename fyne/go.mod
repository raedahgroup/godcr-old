module github.com/raedahgroup/fyne

go 1.12

require (
	fyne.io/fyne v1.1.0
	github.com/decred/slog v1.0.0
	github.com/raedahgroup/dcrlibwallet v1.1.1-0.20190911210329-068758c38d77
	github.com/raedahgroup/godcr/fyne v0.0.0-00010101000000-000000000000
)

replace github.com/raedahgroup/godcr/fyne => ./
