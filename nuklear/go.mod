module github.com/raedahgroup/godcr/nuklear

go 1.12

require (
	github.com/aarzilli/nucular v0.0.0-20190902135428-56f96409e78a
	github.com/atotto/clipboard v0.1.2
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/slog v1.0.0
	github.com/golang/freetype v0.0.0-20161208064710-d9be45aaf745
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190807181808-37b6666fe764
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	golang.org/x/image v0.0.0-20190902063713-cb417be4ba39
	golang.org/x/mobile v0.0.0-20190830201351-c6da95954960
)

replace github.com/raedahgroup/godcr/app => ../app
