module github.com/raedahgroup/godcr/cmd/godcr-cli

go 1.12

require (
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190807181808-37b6666fe764
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/cli v0.0.0-20190912053213-48fdd185f0dd
)

replace github.com/raedahgroup/godcr/app => ../../app
