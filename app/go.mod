module github.com/raedahgroup/godcr/app

go 1.13

replace github.com/raedahgroup/godcr/cli => ../cli

require (
	github.com/decred/dcrd/chaincfg/chainhash v1.0.2
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/hdkeychain v1.1.1
	github.com/decred/dcrd/wire v1.3.0
	github.com/decred/dcrwallet v1.2.3-0.20181120205657-8690f1096aa7
	github.com/decred/dcrwallet/rpc/walletrpc v0.3.0
	github.com/decred/dcrwallet/wallet v1.3.0
	github.com/decred/dcrwallet/walletseed v1.0.3
	github.com/jessevdk/go-flags v1.4.0
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190831020110-aad933e3f96d
	github.com/raedahgroup/godcr/cli v0.0.0-20200107105444-bd23847c1453
	google.golang.org/grpc v1.26.0
)
