module github.com/raedahgroup/godcr/cmd/godcr-nuklear

go 1.13

replace github.com/raedahgroup/godcr/app => ../../app

require (
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/raedahgroup/godcr/app v0.0.0-20200107105444-bd23847c1453
	github.com/raedahgroup/godcr/cli v0.0.0-20200107105444-bd23847c1453
	github.com/raedahgroup/godcr/nuklear v0.0.0-20200107105444-bd23847c1453
)
