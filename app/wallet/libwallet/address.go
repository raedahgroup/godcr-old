package libwallet

import (
	"github.com/raedahgroup/dcrlibwallet"
)

func (lw *LibWallet) AddressInfo(address string) (*dcrlibwallet.AddressInfo, error) {
	return lw.dcrlw.AddressInfo(address)
}

func (lw *LibWallet) ValidateAddress(address string) (bool, error) {
	return lw.dcrlw.IsAddressValid(address), nil
}

func (lw *LibWallet) CurrentReceiveAddress(account uint32) (string, error) {
	return lw.dcrlw.CurrentAddress(int32(account))
}

func (lw *LibWallet) GenerateNewAddress(account uint32) (string, error) {
	return lw.dcrlw.NextAddress(int32(account))
}
