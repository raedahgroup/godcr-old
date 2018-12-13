module github.com/raedahgroup/dcrcli

require (
	github.com/btcsuite/go-flags v0.0.0-20150116065318-6c288d648c1c
	github.com/decred/dcrd/chaincfg/chainhash v1.0.1
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/decred/dcrd/wire v1.2.0
	github.com/decred/dcrwallet v1.2.3-0.20181120205657-8690f1096aa7
	github.com/decred/dcrwallet/rpc/walletrpc v0.1.0
	github.com/decred/dcrwallet/wallet v1.1.0 // indirect
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/raedahgroup/dcrcli/walletrpcclient v0.0.0-20181213135451-898d7ae57860
	github.com/raedahgroup/mobilewallet v0.0.0-20181127040504-952c748fbd60
	github.com/skip2/go-qrcode v0.0.0-20171229120447-cf5f9fa2f0d8
	google.golang.org/genproto v0.0.0-20180928223349-c7e5094acea1 // indirect
	google.golang.org/grpc v1.14.0
)

replace github.com/raedahgroup/mobilewallet => /Users/Apple/go/src/github.com/raedahgroup/mobilewallet
