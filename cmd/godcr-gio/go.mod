module github.com/raedahgroup/godcr/cmd/godcr-gio

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/slog v1.0.0
	github.com/jrick/logrotate v1.0.0
	github.com/raedahgroup/godcr/gio v0.0.0-00010101000000-000000000000
)

replace github.com/raedahgroup/godcr/gio => ../../gio

replace github.com/raedahgroup/dcrlibwallet => github.com/C-ollins/mobilewallet v0.0.0-20191116012520-cf18a67c7aa6

replace github.com/raedahgroup/dcrlibwallet/spv => github.com/C-ollins/mobilewallet/spv v0.0.0-20191116012520-cf18a67c7aa6
