module github.com/raedahgroup/godcr/terminal

go 1.12

require (
	github.com/atotto/clipboard v0.1.2
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/decred/slog v1.0.0
	github.com/gdamore/tcell v1.1.1
	github.com/raedahgroup/dcrlibwallet v1.1.1-0.20190911210329-068758c38d77
	github.com/rivo/tview v0.0.0-20190113120821-e5e361b9d790
)

replace github.com/raedahgroup/dcrlibwallet => ../../dcrlibwallet
