package dcrwalletrpc

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/app"
)

func (c *WalletPRCClient) NetType() string {
	return c.netType
}

func (c *WalletPRCClient) WalletExists() (bool, error) {
	// for now, assume that a wallet has been created since we're connecting through dcrwallet daemon
	// ideally, we'd have to use dcrwallet's WalletLoaderService to do confirm
	return true, nil
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

func (c *WalletPRCClient) IsWalletOpen() bool {
	// for now, assume that the wallet's already open since we're connecting through dcrwallet daemon
	// ideally, we'd have to use dcrwallet's WalletLoaderService to do this
	return true
}

func (c *WalletPRCClient) SyncBlockChain(listener *app.BlockChainSyncListener) error {
	return fmt.Errorf("not yet implemented")
}
