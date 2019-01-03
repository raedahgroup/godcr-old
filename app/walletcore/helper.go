package walletcore

import (
	"errors"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/txscript"
)

func GetAddressFromPkScript(pkScript []byte) (address string, err error) {
	_, addresses, _, err := txscript.ExtractPkScriptAddrs(txscript.DefaultScriptVersion,
		pkScript, currentNet())
	if err != nil {
		return
	}
	if len(addresses) < 1 {
		return "", errors.New("Cannot extract any address from output")
	}
	return addresses[0].EncodeAddress(), nil
}

func currentNet() *chaincfg.Params {
	return &chaincfg.TestNet3Params //Todo check config to see if testnet is active
}
