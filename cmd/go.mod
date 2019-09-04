module github.com/raedahgroup/godcr/cmd

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/raedahgroup/dcrlibwallet v1.1.0 // indirect
	github.com/raedahgroup/godcr v0.0.0-20190904004118-22562be14f04
)

replace github.com/raedahgroup/godcr/fyne => ../fyne
