package gio

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
)

func LoadWallet(appDataDir, netType string) (*dcrlibwallet.MultiWallet, bool, bool, error) {
	multiWallet, err := dcrlibwallet.NewMultiWallet(appDataDir, "", netType)
	if err != nil {
		return nil, false, false, fmt.Errorf("Initialization error: %v", err)
	}

	if multiWallet.LoadedWalletsCount() == 0 {
		return multiWallet, true, false, nil
	}

	var pubPass []byte
	if multiWallet.ReadBoolConfigValueForKey(dcrlibwallet.IsStartupSecuritySetConfigKey, true) {
		// prompt user for public passphrase and assign to `pubPass`
		return multiWallet, false, true, nil
	}

	err = multiWallet.OpenWallets(pubPass)
	if err != nil {
		return multiWallet, false, false, fmt.Errorf("Error opening wallet db: %v", err)
	}

	err = multiWallet.SpvSync()
	if err != nil {
		return multiWallet, false, false, fmt.Errorf("Spv sync attempt failed: %v", err)
	}

	return multiWallet, false, false, nil
}
