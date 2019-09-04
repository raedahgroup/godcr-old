module github.com/raedahgroup/godcr/cmd

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/cli v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/nuklear v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/terminal v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/web v0.0.0-00010101000000-000000000000
)

replace (
	github.com/raedahgroup/godcr/app => ../app
	github.com/raedahgroup/godcr/cli => ../cli
	github.com/raedahgroup/godcr/nuklear => ../nuklear
	github.com/raedahgroup/godcr/terminal => ../terminal
	github.com/raedahgroup/godcr/web => ../web
)
