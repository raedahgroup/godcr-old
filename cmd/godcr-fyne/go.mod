module github.com/raedahgroup/godcr/cmd/godcr-fyne

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/raedahgroup/godcr/fyne v0.0.0-00010101000000-000000000000
)

replace github.com/raedahgroup/godcr/fyne => ../../fyne
