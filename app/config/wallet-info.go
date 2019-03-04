package config

import (
	"encoding/json"
	"fmt"
)

type WalletInfo struct {
	DbDir   string
	Network string
	Source  string
	Default bool
}

func (wallet *WalletInfo) MarshalFlag() (string, error) {
	data, err := json.Marshal(wallet)
	if err == nil {
		return string(data), nil
	}
	return "", err
}

func (wallet *WalletInfo) UnmarshalFlag(value string) error {
	return json.Unmarshal([]byte(value), wallet)
}

func (wallet *WalletInfo) Summary() string {
	return fmt.Sprintf("%s wallet from %s", wallet.Network, wallet.Source)
}

func DefaultWallet(wallets []*WalletInfo) *WalletInfo {
	if len(wallets) == 0 {
		return nil
	}

	for _, wallet := range wallets {
		if wallet.Default {
			return wallet
		}
	}

	return wallets[0]
}
