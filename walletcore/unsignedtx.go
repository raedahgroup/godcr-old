package walletcore

import (
	"errors"
	"fmt"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/txscript"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrwallet/wallet/txrules"
)

func NewUnsignedTx(inputs []*wire.TxIn, sendAmount int64, destinationAddress string,
	changeAddress string) (*wire.MsgTx, error) {

	outputs, err := getTxOutputs(sendAmount, destinationAddress)
	if err != nil {
		return nil, err
	}

	changeSource, err := getChangeSource(changeAddress)
	if err != nil {
		return nil, err
	}

	changeScript, changeScriptVersion, err := changeSource.Script()
	if err != nil {
		return nil, err
	}
	changeScriptSize := changeSource.ScriptSize()
	
	var totalInputAmount int64
	scriptSizes := make([]int, 0, len(inputs))
	for  _, txIn := range(inputs) {
		totalInputAmount += txIn.ValueIn
		scriptSizes = append(scriptSizes, RedeemP2PKHSigScriptSize)
	}

	relayFeePerKb := txrules.DefaultRelayFeePerKb
	maxSignedSize := EstimateSerializeSize(scriptSizes, outputs, changeScriptSize)
	maxRequiredFee := txrules.FeeForSerializeSize(relayFeePerKb, maxSignedSize)
	changeAmount := totalInputAmount - sendAmount - int64(maxRequiredFee)

	if changeAmount != 0 && !txrules.IsDustAmount(dcrutil.Amount(changeAmount), changeScriptSize, relayFeePerKb) {
		if len(changeScript) > txscript.MaxScriptElementSize {
			return nil, errors.New("script size exceed maximum bytes pushable to the stack")
		}
		change := &wire.TxOut{
			Value:    changeAmount,
			Version:  changeScriptVersion,
			PkScript: changeScript,
		}
		outputs = append(outputs, change)
	}

	unsignedTransaction := &wire.MsgTx{
		SerType:  wire.TxSerializeFull,
		Version:  wire.TxVersion, // dcrwallet uses a custom private var txauthor.generatedTxVersion
		TxIn:     inputs,
		TxOut:    outputs,
		LockTime: 0,
		Expiry:   0,
	}

	return unsignedTransaction, nil
}

func getTxOutputs(sendAmount int64, destinationAddress string) ([]*wire.TxOut, error) {
	// get address public script
	pkScript, err := GetPKScript(destinationAddress)
	if err != nil {
		return nil, err
	}

	// create non-change output
	output := wire.NewTxOut(sendAmount, pkScript)
	return []*wire.TxOut{output}, nil
}

func getChangeSource(changeAddress string) (*txChangeSource, error) {
	// get address public script
	pkScript, err := GetPKScript(changeAddress)
	if err != nil {
		return nil, err
	}

	// generate address in source account to receive change
	return &txChangeSource{
		script:  pkScript,
		version: txscript.DefaultScriptVersion,
	}, nil
}

func GetPKScript(address string) ([]byte, error) {
	// decode change address
	addr, err := dcrutil.DecodeAddress(address)
	if err != nil {
		return nil, fmt.Errorf("error decoding change address: %s", err.Error())
	}
	// get change address script
	return txscript.PayToAddrScript(addr)
}
