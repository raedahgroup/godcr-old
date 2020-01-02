module github.com/raedahgroup/godcr/gio

go 1.12

require (
	//gioui.org v0.0.0-20191023081811-143d2aae95dd
	gioui.org v0.0.0-20191211234536-7814da47a0ff
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/hdkeychain v1.1.1
	github.com/decred/dcrwallet/walletseed v1.0.1
	github.com/decred/slog v1.0.0
	github.com/raedahgroup/dcrlibwallet v0.0.0-00010101000000-000000000000
	//github.com/raedahgroup/dcrlibwallet v1.1.1-0.20191224211651-29f1df229b2a
	golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
)

//replace github.com/raedahgroup/dcrlibwallet => github.com/C-ollins/mobilewallet v0.0.0-20191116012520-cf18a67c7aa6
eplace github.com/raedahgroup/dcrlibwallet => github.com/C-ollins/mobilewallet v1.0.0-rc1.0.2019206032901-ef455a3cc250
replace github.com/raedahgroup/dcrlibwallet/spv => github.com/C-ollins/mobilewallet/spv v0.0.0-2019206032901-ef455a3cc250