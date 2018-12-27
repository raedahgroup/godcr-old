package config

func configText() string {
	return `[Application Options]

; ------------------------------------------------------------------------------
; Network settings
; ------------------------------------------------------------------------------

; Wallet RPC server
walletrpcserver={{.RpcAddress}}

; Username and password to authenticate connections to a Dcrwallet RPC server
rpcuser={{.RpcUsername}}
rpcpass={{.RpcPassword}}

; Specify dcrwallet's RPC certificate, or disable TLS for the connection
; rpccert=/home/me/.dcrwallet/rpc.cert
; nodaemontls=0

; ------------------------------------------------------------------------------
; HTTP Server
; ------------------------------------------------------------------------------

; Server address to serve dcrcli via HTTP (when http=true or http=1)
serveraddress={{.ServerAddress}}

; Turn HTTP mode on (http=true or http=1) or off (http=false or http=0).
; HTTP mode is off by default.
; http=true`

}
