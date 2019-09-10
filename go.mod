module github.com/raedahgroup/godcr

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/decred/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/raedahgroup/godcr/app v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/cli v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/fyne v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/nuklear v0.0.0-20190904233416-f0084ce3a199
	github.com/raedahgroup/godcr/terminal v0.0.0-20190904233416-f0084ce3a199
	github.com/raedahgroup/godcr/web v0.0.0-20190904233416-f0084ce3a199
)

replace (
	github.com/raedahgroup/godcr/app => ./app
	github.com/raedahgroup/godcr/cli => ./cli
	github.com/raedahgroup/godcr/fyne => ./fyne
)
