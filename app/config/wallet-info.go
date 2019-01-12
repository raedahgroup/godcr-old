package config

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
)

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
	// load default config values and create parser object with it
	config := defaultConfig()
	parser := flags.NewParser(&config, flags.None)

	// read current config file content into config object
	fileParser := flags.NewIniParser(parser)
	err = fileParser.ParseFile("./godcr.conf")
	if err != nil {
		return fmt.Errorf("Error reading config file: %s", err.Error())
	}

	// remove previous wallets info
	printConfigGroup(parser.Groups())

	// add new wallets info
	for _, w := range wallets {
		parser.AddGroup("Detected Wallet", w.LongDescription(), w)
	}

	// write config object to file
	err = fileParser.WriteFile("./godcr.conf", flags.IniIncludeComments|flags.IniIncludeDefaults|flags.IniCommentDefaults)
	if err != nil {
		return fmt.Errorf("Error saving changes to config file: %s", err.Error())
	}

	return
}

func printConfigGroup(groups []*flags.Group) {
	for _, configGroup := range groups {
		fmt.Println(configGroup.ShortDescription, len(configGroup.Options()), "options")

		for _, option := range configGroup.Options() {
			fmt.Println(option.String(), option.Value())
		}

		if len(configGroup.Groups()) > 0 {
			printConfigGroup(configGroup.Groups())
		}
	}
}
