package helper 

import (
	"fmt"
	"github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrwallet/walletseed"
)

func GenerateSeedWords() (string, error) {
	// generate seed
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return "", fmt.Errorf("\nError generating seed for new wallet: %s.", err)
	}
	return walletseed.EncodeMnemonic(seed), nil
}