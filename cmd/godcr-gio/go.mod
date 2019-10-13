module github.com/raedahgroup/godcr/cmd/godcr-gio

go 1.12

require (
	github.com/raedahgroup/godcr/app v0.0.0-20191001132534-0f6d1a0712a5
	github.com/raedahgroup/godcr/cli v0.0.0-00010101000000-000000000000
	github.com/raedahgroup/godcr/gio v0.0.0-00010101000000-000000000000
)

replace github.com/raedahgroup/godcr/gio => ../../gio

replace github.com/raedahgroup/godcr/cli => ../../cli

replace github.com/raedahgroup/godcr/app => ../../app
