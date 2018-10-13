module github.com/codemaestro64/dcrcli

require (
	github.com/btcsuite/go-flags v0.0.0-20150116065318-6c288d648c1c
	github.com/decred/dcrd/dcrutil v1.1.1
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/raedahgroup/dcrcli/walletrpcclient v0.0.1
	golang.org/x/sys v0.0.0-20180928133829-e4b3c5e90611 // indirect
	google.golang.org/genproto v0.0.0-20180928223349-c7e5094acea1 // indirect
)

replace github.com/raedahgroup/dcrcli/walletrpcclient => ./walletrpcclient
