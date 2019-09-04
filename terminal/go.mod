module github.com/raedahgroup/godcr/terminal

go 1.12

require (
	github.com/atotto/clipboard v0.1.2
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/gdamore/tcell v1.1.1
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190807181808-37b6666fe764
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/rivo/tview v0.0.0-20190113120821-e5e361b9d790
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
)

replace github.com/raedahgroup/godcr/app => ../app
