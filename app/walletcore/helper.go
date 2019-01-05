package walletcore

import (
	"errors"
	"strings"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/txscript"
)

func GetAddressFromPkScript(activeNet *chaincfg.Params, pkScript []byte) (address string, err error) {
	_, addresses, _, err := txscript.ExtractPkScriptAddrs(txscript.DefaultScriptVersion,
		pkScript, activeNet)
	if err != nil {
		return
	}
	if len(addresses) < 1 {
		return "", errors.New("Cannot extract any address from output")
	}

	encodedAddresses := make([]string, len(addresses))
	for i, address := range addresses {
		encodedAddresses[i] = address.EncodeAddress()
	}

	return strings.Join(encodedAddresses, ", "), nil
}
