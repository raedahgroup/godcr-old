package dcrlibwallet

func (lib *DcrWalletLib) SaveToSettings(key string, value interface{}) error {
	return lib.walletLib.SaveToSettings(key, value)
}

func (lib *DcrWalletLib) ReadFromSettings(key string, valueOut interface{}) (error) {
	return lib.walletLib.ReadFromSettings(key, valueOut)
}
