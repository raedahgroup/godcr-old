module github.com/raedahgroup/dcrcli

require (
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/jessevdk/go-flags v1.4.0
	github.com/raedahgroup/dcrcli/walletrpcclient v0.0.1
	github.com/skip2/go-qrcode v0.0.0-20171229120447-cf5f9fa2f0d8
)

replace github.com/raedahgroup/dcrcli/walletrpcclient => ./walletrpcclient
