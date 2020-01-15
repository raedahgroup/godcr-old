module github.com/raedahgroup/godcr/cmd/godcr-gio

go 1.12

require (
	//gioui.org v0.0.0-20200113204813-a7dc7c01c0f5 // indirect
	gioui.org v0.0.0-20191211234536-7814da47a0ff
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/slog v1.0.0
	github.com/gobuffalo/envy v1.8.1 // indirect
	github.com/gobuffalo/logger v1.0.3 // indirect
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/gobuffalo/packr/v2 v2.7.1 // indirect
	github.com/jrick/logrotate v1.0.0
	github.com/karrick/godirwalk v1.14.0 // indirect
	github.com/raedahgroup/dcrlibwallet v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/gio v0.0.0-00010101000000-000000000000
	github.com/rogpeppe/go-internal v1.5.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.0.0-20200109152110-61a87790db17 // indirect
	golang.org/x/sys v0.0.0-20200113162924-86b910548bc1 // indirect
)

replace github.com/raedahgroup/godcr/gio => ../../gio

replace github.com/raedahgroup/dcrlibwallet => github.com/C-ollins/mobilewallet v1.0.0-rc1.0.20191206032901-ef455a3cc250

replace github.com/raedahgroup/dcrlibwallet/spv => github.com/C-ollins/mobilewallet/spv v0.0.0-20191206032901-ef455a3cc250
