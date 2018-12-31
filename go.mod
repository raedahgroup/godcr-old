module github.com/raedahgroup/godcr

require (
	github.com/aarzilli/nucular v0.0.0-20181227101716-d1a942545d6d
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/jessevdk/go-flags v1.4.0
	github.com/raedahgroup/godcr/walletrpcclient v0.0.1
	github.com/skip2/go-qrcode v0.0.0-20171229120447-cf5f9fa2f0d8
	golang.org/x/image v0.0.0-20181116024801-cd38e8056d9b
)

replace github.com/raedahgroup/godcr/walletrpcclient => ./walletrpcclient
