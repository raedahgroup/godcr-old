module github.com/raedahgroup/godcr/cmd/godcr-fyne

go 1.12

replace github.com/raedahgroup/godcr/fyne => ../../fyne

replace github.com/raedahgroup/dcrlibwallet => github.com/C-ollins/mobilewallet v1.0.0-rc1.0.20191206032901-ef455a3cc250

replace github.com/raedahgroup/dcrlibwallet/spv => github.com/C-ollins/mobilewallet/spv v0.0.0-20191206032901-ef455a3cc250

require (
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/raedahgroup/godcr/fyne v0.0.0-00010101000000-000000000000
)
