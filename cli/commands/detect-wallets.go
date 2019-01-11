package commands

import (
	"fmt"
	"github.com/raedahgroup/godcr/app/config"
	"os"
	"path/filepath"
)

type DetectWalletsCommand struct {}

func (detectCmd DetectWalletsCommand) Execute(args []string) error {
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
