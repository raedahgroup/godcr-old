package dcrwalletrpc

import (
	"context"
	"fmt"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrcli/app"
)

func (c *WalletPRCClient) NetType() string {
	return c.netType
}

func (c *WalletPRCClient) WalletExists() (bool, error) {
	res, err := c.walletLoader.WalletExists(context.Background(), &walletrpc.WalletExistsRequest{})
	if err != nil {
		return false, err
	}
	return res.Exists, nil
}

func (c *WalletPRCClient) GenerateNewWalletSeed() (string, error) {
	return "", fmt.Errorf("not yet implemented")
}

func (c *WalletPRCClient) CreateWallet(passphrase, seed string) error {
	// since we're connecting through dcrwallet daemon, assume that the wallet's already been created
	// calling create again should return an error
	// ideally, we'd have to use dcrwallet's WalletLoaderService to do this
	return fmt.Errorf("wallet should already be created by dcrwallet daemon")
}

func (c *WalletPRCClient) OpenWallet() error {
	// for now, assume that the wallet's already open since we're connecting through dcrwallet daemon
	// ideally, we'd have to use dcrwallet's WalletLoaderService to do this
	return nil
}

func (c *WalletPRCClient) CloseWallet() {}

func (c *WalletPRCClient) IsWalletOpen() bool {
	// for now, assume that the wallet's already open since we're connecting through dcrwallet daemon
	// ideally, we'd have to use dcrwallet's WalletLoaderService to do this
	return true
}

func (c *WalletPRCClient) SyncBlockChain(listener *app.BlockChainSyncListener, showLog bool) error {
	// pretend to start and successfully complete sync
	listener.SyncStarted()
	listener.SyncEnded(nil)
	return nil
}
