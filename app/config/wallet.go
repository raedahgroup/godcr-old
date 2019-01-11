package config

import "fmt"

type WalletInfo struct {
	DbPath string
	NetType string
	Source string
}

func (w WalletInfo) String() string {
	return fmt.Sprintf("Path: %s\nType: %s\nSource: %s\n", w.DbPath, w.NetType, w.Source)
}
