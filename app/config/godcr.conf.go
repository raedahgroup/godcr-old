package config

const configTextTemplate = `[Application Options]

; ------------------------------------------------------------------------------
; App Data
; ------------------------------------------------------------------------------

; Path to application data directory.
appdata={{.AppDataDir}}

; ------------------------------------------------------------------------------
; Network settings
; ------------------------------------------------------------------------------

; Connect to a running drcwallet daemon over rpc to perform wallet operations.
; By default godcr uses dcrlibwallet (usewalletrpc=false or usewalletrpc=0).
; You can switch to dcrwallet rpc using (usewalletrpc=true or usewalletrpc=1).
; usewalletrpc=false

; Wallet gRPC server address. Required if usewalletrpc=true
; walletrpcserver=localhost:19111

; Disable TLS for rpc connection (nowalletrpctls=1) or specify path to dcrwallet's RPC certificate
; nowalletrpctls=0
walletrpccert={{.WalletRPCCert}}

; Connects to testnet wallet instead of mainnet
; testnet=false

; ------------------------------------------------------------------------------
; Godcr Interface Modes
; ------------------------------------------------------------------------------

; godcr can run in any of the following modes: cli, http, desktop
; mode=cli

; ------------------------------------------------------------------------------
; Godcr HTTP Mode Options
; ------------------------------------------------------------------------------

; Host and port for godcr http server when mode=http
httphost={{.HTTPHost}}
httpport={{.HTTPPort}}`