module github.com/raedahgroup/godcr/nuklear

go 1.12

require (
	github.com/aarzilli/nucular v0.0.0-20190403084742-0071461892e4
	github.com/atotto/clipboard v0.1.2
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/decred/slog v1.0.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190807181808-37b6666fe764
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	golang.org/x/image v0.0.0-20190501045829-6d32002ffd75
	golang.org/x/mobile v0.0.0-20190318164015-6bd122906c08
)

replace github.com/raedahgroup/godcr/app => ../app
