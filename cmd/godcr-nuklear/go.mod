module github.com/raedahgroup/godcr/cmd/godcr-nuklear

go 1.12

require (
	github.com/decred/dcrd/dcrutil v1.2.0
	github.com/gobuffalo/packr/v2 v2.7.0
	github.com/raedahgroup/godcr/nuklear v0.0.0-00010101000000-000000000000
)

replace github.com/raedahgroup/godcr/nuklear => ../../nuklear
