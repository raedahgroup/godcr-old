module github.com/raedahgroup/godcr/cli

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/hdkeychain v1.1.1
	github.com/decred/dcrwallet v1.2.3-0.20181120205657-8690f1096aa7
	github.com/decred/dcrwallet/walletseed v1.0.3
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190831020110-aad933e3f96d
	github.com/raedahgroup/godcr/app v0.0.0-20200107105444-bd23847c1453
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
)

replace github.com/raedahgroup/godcr/app => ../app
