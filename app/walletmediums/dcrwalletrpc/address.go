package dcrwalletrpc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app/config"
)

type ExplicitString struct {
	Value         string
	explicitlySet bool
}

type walletConfig struct {
	GRPCListeners      []string `long:"grpclisten" description:"Listen for gRPC connections on this interface/port"`
	LegacyRPCListeners []string `long:"rpclisten" description:"Listen for legacy JSON-RPC connections on this interface/port"`
	NoGRPC             bool     `long:"nogrpc" description:"Disable the gRPC server"`
	NoLegacyRPC        bool     `long:"nolegacyrpc" description:"Disable the legacy JSON-RPC server"`
	DisableServerTLS       bool                    `long:"noservertls" description:"Disable TLS for the RPC servers -- NOTE: This is only allowed if the RPC server is bound to localhost"`
	RPCCert                *ExplicitString `long:"rpccert" description:"File containing the certificate file"`

	TBOpts ticketBuyerOptions `group:"Ticket Buyer Options" namespace:"ticketbuyer"`
}

type ticketBuyerOptions struct{}

const (
	defaultWalletConfigFilename = "dcrwallet.conf"
)

var (
	walletConfigFilePath = filepath.Join(config.DefaultDcrwalletAppDataDir, defaultWalletConfigFilename)
)

func walletAddressFromDcrdwalletConfig() (addresses []string, err error) {
	wConfig := walletConfig{}

	parser := flags.NewParser(&wConfig, flags.IgnoreUnknown)
	err = flags.NewIniParser(parser).ParseFile(walletConfigFilePath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			err = fmt.Errorf("Error parsing configuration file: %v", err.Error())
			return
		}
		return
	}

	if !wConfig.NoGRPC {
		return wConfig.GRPCListeners, nil
	}

	return wConfig.LegacyRPCListeners, nil
}
