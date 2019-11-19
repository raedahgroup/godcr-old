module github.com/raedahgroup/godcr/fyne

go 1.13

require (
	fyne.io/fyne v1.1.3-0.20191104221827-e8f6795efa08
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/slog v1.0.0
	github.com/gen2brain/beeep v0.0.0-20190719094215-ece0cb67ca77
	github.com/go-gl/glfw v0.0.0-20190409004039-e6da0acd62b1 // indirect
	github.com/gobuffalo/packr/v2 v2.6.0
	github.com/raedahgroup/dcrlibwallet v0.0.0-00010101000000-000000000000
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
)

replace github.com/raedahgroup/dcrlibwallet => github.com/C-ollins/mobilewallet v0.0.0-20191116012520-cf18a67c7aa6

replace github.com/raedahgroup/dcrlibwallet/spv => github.com/C-ollins/mobilewallet/spv v0.0.0-20191116012520-cf18a67c7aa6
