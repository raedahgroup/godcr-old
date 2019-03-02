package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

var walletConfigFilePath = filepath.Join(DefaultAppDataDir, "wallets.conf")

type WalletInfo struct {
	DbPath string `long:"db"`
	NetType string `long:"nettype"`
	Source string `long:"source"`
}

func (w WalletInfo) ShortDescription(index int) string {
	return fmt.Sprintf("wallet_%d", index)
}

func (w WalletInfo) LongDescription() string {
	return fmt.Sprintf("%s %s wallet", w.Source, w.NetType)
}

func SaveDetectedWalletsInfo(wallets []*WalletInfo) (err error) {
	viper.SetConfigFile("./godcr.conf")
	err = fileParser.ParseFile()
	if err != nil {
		return fmt.Errorf("Error reading config file: %s", err.Error())
	}

	// add new wallets info
	for i, w := range wallets {
		.AddGroup(fmt.Sprintf("Detected Wallet %d", i+1), w.LongDescription(), w)
	}

	// write config object to file
	err = fileParser.WriteFile("./godcr.conf", flags.IniIncludeComments|flags.IniIncludeDefaults|flags.IniCommentDefaults)
	if err != nil {
		return fmt.Errorf("Error saving changes to config file: %s", err.Error())
	}

	return
}
