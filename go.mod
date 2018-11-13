module github.com/raedahgroup/dcrcli

require (
	github.com/Baozisoftware/qrcode-terminal-go v0.0.0-20170407111555-c0650d8dff0f
	github.com/btcsuite/go-flags v0.0.0-20150116065318-6c288d648c1c
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/decred/dcrd/dcrutil v1.1.1
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/juju/ansiterm v0.0.0-20180109212912-720a0952cc2a // indirect
	github.com/lunixbochs/vtclean v0.0.0-20180621232353-2d01aacdc34a // indirect
	github.com/manifoldco/promptui v0.3.1
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/raedahgroup/dcrcli/walletrpcclient v0.0.1
	github.com/skip2/go-qrcode v0.0.0-20171229120447-cf5f9fa2f0d8
	golang.org/x/sys v0.0.0-20180928133829-e4b3c5e90611 // indirect
	google.golang.org/genproto v0.0.0-20180928223349-c7e5094acea1 // indirect
	google.golang.org/grpc v1.14.0
)

replace github.com/raedahgroup/dcrcli/walletrpcclient => ./walletrpcclient
