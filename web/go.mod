module github.com/raedahgroup/godcr/web

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/decred/dcrwallet/wallet v1.3.0
	github.com/decred/slog v1.0.0
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gobuffalo/logger v1.0.1 // indirect
	github.com/gobuffalo/packr v1.30.1
	github.com/gobuffalo/packr/v2 v2.6.0
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190807181808-37b6666fe764
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/cli v0.0.0-00010101000000-000000000000
	github.com/rogpeppe/go-internal v1.3.1 // indirect
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	github.com/stretchr/testify v1.4.0 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/sys v0.0.0-20190904154756-749cb33beabd // indirect
	google.golang.org/appengine v1.6.2 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	github.com/raedahgroup/godcr/app => ../app
	github.com/raedahgroup/godcr/cli => ../cli
)
