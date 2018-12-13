package mobilewallet

import (
	"errors"
	"github.com/decred/dcrwallet/wallet"

	ws "github.com/raedahgroup/dcrcli/wallet"
)

type Client struct {
	
}

// Balance implements `walletsource.Balance`
func (c *Client) Balance(accountNumbers ...uint32) ([]*wallet.AccountBalance, error) {
	return nil, errors.New("Not yet implemented")
}

// NextAccount implements `walletsource.NextAccount`
func (c *Client) NextAccount(accountName string, passphrase string) (uint32, error) {
	return 0, errors.New("Not yet implemented")
}

// GenerateReceiveAddress implements `walletsource.GenerateReceiveAddress`
func (c *Client) GenerateReceiveAddress(account uint32) (string, error) {
	return "", errors.New("Not yet implemented")
}

// ValidateAddress implements `walletsource.ValidateAddress`
func (c *Client) ValidateAddress(address string) (bool, error) {
	return false, errors.New("Not yet implemented")
}

// AccountSend implements `walletsource.NextAccount`
func (c *Client) AccountSend(amountInDCR float64, sourceAccount uint32, destinationAddress, passphrase string) (string, error) {
	return "", errors.New("Not yet implemented")
}

// UnspentOutputs implements `walletsource.NextAccount`
func (c *Client) UnspentOutputs(account uint32, targetAmount int64) ([]*wallet.UnspentOutput, error) {
	return nil, errors.New("Not yet implemented")
}

// Send implements `walletsource.NextAccount`
func (c *Client) UTXOSend(utxoKeys []string, amountInDCR float64, sourceAccount uint32, destinationAddress, passphrase string) (string, error) {
	return "", errors.New("Not yet implemented")
}