package helper 

import (
	"fmt"
	"github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrwallet/walletseed"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
)

type (
	MultiWallet struct {
		*dcrlibwallet.MultiWallet 
		walletIDs []int
	}
)

func LoadWallet(appDataDir, netType string) (*MultiWallet, bool, bool, error) {
	multiWallet, err := dcrlibwallet.NewMultiWallet(appDataDir, "", netType)
	if err != nil {
		return nil, false, false, fmt.Errorf("Initialization error: %v", err)
	}
	
	mw := &MultiWallet{
		MultiWallet: multiWallet,
	} 
	mw.walletIDs = make([]int, 0)
	mw.walletIDs = append(mw.walletIDs, mw.OpenedWalletIDsRaw()...)

	if multiWallet.LoadedWalletsCount() == 0 {
		return mw, true, false, nil
	}

	var pubPass []byte
	if multiWallet.ReadBoolConfigValueForKey(dcrlibwallet.IsStartupSecuritySetConfigKey, true) {
		// prompt user for public passphrase and assign to `pubPass`
		//return mw, false, true, nil
	}

	err = multiWallet.OpenWallets(pubPass)
	if err != nil {
		return mw, false, false, fmt.Errorf("Error opening wallet db: %v", err)
	}

	for i := range mw.walletIDs {
		fmt.Println(mw.WalletWithID(mw.walletIDs[i]).WalletOpened())
	}

	err = multiWallet.SpvSync()
	if err != nil {
		return mw, false, false, fmt.Errorf("Spv sync attempt failed: %v", err)
	}

	


	return mw, false, false, nil
}

func (w *MultiWallet) RegisterWalletID(wID int) {
	for _,v := range w.walletIDs {
		if v == wID {
			return
		}
	}
	
	w.walletIDs = append(w.walletIDs, wID)
	// TODO return and handle wallet is already registered error
}

func (w *MultiWallet) TotalBalance() (string, error) {
	var totalBalance int64 

	for _, walletID := range w.walletIDs {
		accounts, err := w.WalletWithID(walletID).GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			return "0", err
		}

		for _, account := range accounts.Acc {
			totalBalance += account.TotalBalance
		}
	}

	return dcrutil.Amount(totalBalance).String(), nil
}

func GenerateSeedWords() (string, error) {
	// generate seed
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return "", fmt.Errorf("\nError generating seed for new wallet: %s.", err)
	}
	return walletseed.EncodeMnemonic(seed), nil
}