module github.com/raedahgroup/godcr

require (
	github.com/aarzilli/nucular v0.0.0-20190403084742-0071461892e4
	github.com/atotto/clipboard v0.1.2
	github.com/decred/dcrd/chaincfg/chainhash v1.0.1
	github.com/decred/dcrd/dcrjson v1.2.0 // indirect
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/dcrd/hdkeychain v1.1.1
	github.com/decred/dcrd/wire v1.2.0
	github.com/decred/dcrwallet v1.2.3-0.20181120205657-8690f1096aa7
	github.com/decred/dcrwallet/rpc/walletrpc v1.0.1-0.20181109211527-ca582da21c08
	github.com/decred/dcrwallet/wallet v1.3.0
	github.com/decred/dcrwallet/walletseed v1.0.1
	github.com/decred/slog v1.0.0
	github.com/gdamore/tcell v1.1.1
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/gorilla/websocket v1.2.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/raedahgroup/dcrlibwallet v1.1.0
	github.com/raedahgroup/godcr/cmd v0.0.0-00010101000000-000000000000 // indirect
	github.com/rivo/tview v0.0.0-20190113120821-e5e361b9d790
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	golang.org/x/image v0.0.0-20190501045829-6d32002ffd75
	golang.org/x/mobile v0.0.0-20190318164015-6bd122906c08
	google.golang.org/grpc v1.19.0
)

replace (
	github.com/raedahgroup/godcr/cmd => ./cmd
	github.com/raedahgroup/godcr/fyne => ./fyne
)
