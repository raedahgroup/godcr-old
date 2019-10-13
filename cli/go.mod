module github.com/raedahgroup/godcr/cli

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/decred/dcrd/hdkeychain v1.1.1
	github.com/decred/dcrwallet v1.2.3-0.20181120205657-8690f1096aa7
	//github.com/decred/dcrwallet/rpc/walletrpc v1.0.1-0.20181109211527-ca582da21c08
	github.com/decred/dcrwallet/rpc/walletrpc v0.2.1-0.20191007153235-6a27c792bbbb // indirect
	github.com/decred/dcrwallet/walletseed v1.0.1
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190807181808-37b6666fe764
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
)

replace github.com/raedahgroup/godcr/app => ../app
