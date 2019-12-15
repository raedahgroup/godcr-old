module github.com/raedahgroup/godcr/fyne

go 1.13

require (
	fyne.io/fyne v1.2.1-0.20191214234542-7e09093e0f38
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/slog v1.0.0
	github.com/gen2brain/beeep v0.0.0-20190719094215-ece0cb67ca77
	github.com/gobuffalo/packr/v2 v2.6.0
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/gopherjs/gopherwasm v1.1.0 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/raedahgroup/dcrlibwallet v0.0.0-00010101000000-000000000000
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	github.com/tadvi/systray v0.0.0-20190226123456-11a2b8fa57af // indirect
	gopkg.in/toast.v1 v1.0.0-20180812000517-0a84660828b2 // indirect
)

replace github.com/raedahgroup/dcrlibwallet => github.com/C-ollins/mobilewallet v1.0.0-rc1.0.20191201141735-f45887d0465f

replace github.com/raedahgroup/dcrlibwallet/spv => github.com/C-ollins/mobilewallet/spv v0.0.0-20191201141735-f45887d0465f
