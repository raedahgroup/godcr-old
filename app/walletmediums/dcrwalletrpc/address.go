package dcrwalletrpc

import (
	"path/filepath"

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

	TBOpts ticketBuyerOptions `group:"Ticket Buyer Options" namespace:"ticketbuyer"`
}

type ticketBuyerOptions struct{}

const (
	defaultWalletConfigFilename = "dcrwallet.conf"
)

var (
	walletConfigFilePath = filepath.Join(config.DefaultDcrwalletAppDataDir, defaultWalletConfigFilename)
)

func walletAddressFromDcrdwalletConfig() (addresses []string, notls bool, certpath string, err error) {
	wConfig := walletConfig{}

	parser := flags.NewParser(&wConfig, flags.IgnoreUnknown)
	err = flags.NewIniParser(parser).ParseFile(walletConfigFilePath)
	if err != nil {
		return
	}
	var rpcCert string
	if wConfig.RPCCert != nil {
		rpcCert = wConfig.RPCCert.Value
	}
	return wConfig.GRPCListeners, wConfig.DisableServerTLS, rpcCert, nil
}
