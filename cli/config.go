package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/decred/dcrd/dcrutil"
)

const (
	defaultConfigFilename    = "dcrcli.conf"
	defaultLogDirname        = "log"
	defaultLogFilename       = "dcrcli.log"
	defaultHTTPServerAddress = "127.0.0.1:1234"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("dcrcli", false)
	defaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(defaultDcrwalletAppDataDir, "rpc.cert")
	defaultConfigFile          = filepath.Join(defaultAppDataDir, defaultConfigFilename)
	defaultLogDir              = filepath.Join(defaultAppDataDir, defaultLogDirname)
)

type Config struct {
	ShowVersion       bool   `short:"v" long:"version" description:"Display version information and exit"`
	ConfigFile        string `short:"C" long:"configfile" description:"Path to configuration file"`
	RPCUser           string `short:"u" long:"rpcuser" description:"RPC username"`
	RPCPassword       string `short:"p" long:"rpcpass" default-mask:"-" description:"RPC password"`
	WalletRPCServer   string `short:"w" long:"walletrpcserver" description:"Wallet RPC server to connect to"`
	RPCCert           string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	HTTPServerAddress string `short:"s" long:"serveraddress" description:"Address and port of the HTTP server."`
	HTTPMode          bool   `long:"http" description:"Run in HTTP mode."`
	NoDaemonTLS       bool   `long:"nodaemontls" description:"Disable TLS"`
}

func DefaultConfig() Config {
	return Config{
		ConfigFile:        defaultConfigFile,
		RPCCert:           defaultRPCCertFile,
		HTTPServerAddress: defaultHTTPServerAddress,
	}
}

func AppName() string {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	return appName
}
