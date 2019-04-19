package dcrwalletrpc

import (
	"path/filepath"

	"context"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app/config"
)

type ExplicitString struct {
	Value         string
	explicitlySet bool
}

type walletConfig struct {
	GRPCListeners    []string        `long:"grpclisten" description:"Listen for gRPC connections on this interface/port"`
	DisableServerTLS bool            `long:"noservertls" description:"Disable TLS for the RPC servers -- NOTE: This is only allowed if the RPC server is bound to localhost"`
	RPCCert          *ExplicitString `long:"rpccert" description:"File containing the certificate file"`

	TBOpts struct{} `group:"Ticket Buyer Options" namespace:"ticketbuyer"`
}

const (
	defaultWalletConfigFilename = "dcrwallet.conf"
)

var (
	walletConfigFilePath = filepath.Join(config.DefaultDcrwalletAppDataDir, defaultWalletConfigFilename)
)

func parseDcrWalletConfigAndConnect(ctx context.Context, useRpcCert string, disableTls bool) *WalletRPCClient {
	wConfig := walletConfig{}

	parser := flags.NewParser(&wConfig, flags.IgnoreUnknown)
	err := flags.NewIniParser(parser).ParseFile(walletConfigFilePath)
	if err != nil {
		return nil
	}

	var rpcCert string
	if useRpcCert != "" {
		rpcCert = useRpcCert
	} else if wConfig.RPCCert != nil {
		rpcCert = wConfig.RPCCert.Value
	}

	for _, address := range wConfig.GRPCListeners {
		walletRPCClient, _ := createConnection(ctx, address, rpcCert, disableTls)
		if walletRPCClient != nil {
			config.UpdateConfigFile(func(config *config.ConfFileOptions) {
				config.WalletRPCServer = address
			})
			return walletRPCClient
		}
	}

	return nil
}
