module github.com/raedahgroup/godcr/web

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/decred/dcrwallet/wallet v1.3.0
	github.com/decred/slog v1.0.0
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gorilla/websocket v1.2.0
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190807181808-37b6666fe764
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/cli v0.0.0-00010101000000-000000000000
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
)

replace (
	github.com/raedahgroup/godcr/app => ../app
	github.com/raedahgroup/godcr/cli => ../cli
)
