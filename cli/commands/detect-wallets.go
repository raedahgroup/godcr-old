package commands

import (
	"context"
	"fmt"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"os"
	"path/filepath"
)

type DetectWalletsCommand struct {
	commanderStub
}

func (detectCmd DetectWalletsCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	err := filepath.Walk("/usr/lib", func(path string, file os.FileInfo, err error) error {
		if err == nil && !file.IsDir() && file.Name() == config.WalletDbFileName {
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
