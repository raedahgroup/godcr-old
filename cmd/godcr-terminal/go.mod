module github.com/raedahgroup/godcr/cmd/godcr-terminal

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/raedahgroup/godcr/terminal v0.0.0-00010101000000-000000000000
)

replace github.com/raedahgroup/godcr/terminal => ../../terminal

replace github.com/raedahgroup/dcrlibwallet => ../../../dcrlibwallet
