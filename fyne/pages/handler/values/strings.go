package values

const (
	//General

	DCR   = "DCR"
	Bytes = "byte"

	// Send Page

	// General
	DestinationAddressPlaceHolder = "Destination address"
	InvalidAddress                = "Invalid address"
	NilAmount                     = "- DCR"
	Send                          = "Send"
	SwitchToSendToAccount         = "Send to self"
	SwitchToSendToAddress         = "Send to others"
	SuccessText                   = "Transaction sent"
	ZeroByte                      = "0 bytes"
	ZeroAmount                    = "0 DCR"
	// Base Object
	GotIt      = "Got it"
	ClearField = "Clear all fields"
	BaseWidget = ""
	// Account Selector
	Imported = "imported"
	// From Account Selector
	FromText                            = "From"
	FromAccountSelectorPopUpHeaderLabel = "Sending account"
	SpendableAmountLabel                = "Spendable: "
	Spendable                           = "Spendable"
	// Amount Entry
	Amount                = "Amount"
	AmountRegExp          = "^\\d*\\.?\\d*$"
	MaxAmountAllowedInDCR = "12345678.12345678"
	Max                   = "MAX"
	NoFunds               = "Not enough funds"
	NoFundsOrNotConnected = "Not enough funds (or not connected)."
	SendPageInfo          = "Input the destination \nwallet address and the amount in \nDCR to send funds."
	SendDcr               = "Send DCR"
	TestAddress           = "HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH" // max width for a 35 alphabeth
	// Confirmation Window
	BalanceAfterSend    = "Balance after send"
	ConfirmToSend       = "Confirm to send"
	FailedToSend        = "Failed to send. Please try again"
	SendingFrom         = "Sending from"
	SendingDcrWarning   = "Your DCR will be sent and CANNOT be undone"
	ToDesinationAddress = "To destination address"
	ToSelf              = "To self"
	TransactionFee      = "Transaction fee"
	TotalCost           = "Total cost"
	// Transaction Info Window
	FeeRate            = "Fee rate"
	FeeRateInfo        = "0.0001 DCR/byte"
	ProcessingTime     = "Processing time"
	ProcessingTimeInfo = "Approx. 10 mins (2 blocks)"
	TransactionSize    = "Transaction size"

	// Password Popup
	Cancel               = "Cancel"
	Confirm              = "Confirm"
	SpendingPasswordText = "Spending Password"

	// Receive Page
	ReceivePageLabel      = "Receive DCR"
	ReceivingAccountLabel = "Receiving account"
	TapToCopy             = "(Tap to copy)"
	GenerateNewAddress    = "Generate new address"
	ReceivePageInfo       = "Each time you request a \npayment, a new address is \ncreated to protect your privacy."

	// Error Messages
	AmountDecimalPlaceErr       = "Amount has more than 8 decimal places"
	AccountSelectorIconErr      = "Could not retrieve account selector icons"
	AccountDetailsErr           = "Could not retrieve account details"
	AccountBalanceErr           = "Could not retrieve account balance"
	AccountNumberErr            = "Could not retrieve account number"
	BaseObjectsIconErr          = "Could not retrieve base object icons"
	ConfirmationWindowIconsErr  = "Could not retrieve confirmation window icons"
	ConfirmationWindowErr       = "Could not view confirmation window"
	GettingAccountBalanceErr    = "Could not retrieve account balace for send destination"
	GettingAddress              = "Could not generate address"
	GettingAddressToSelfSendErr = "Could not generate address to send to self"
	InsufficientBalanceErr      = "Insufficient balance"
	InitTxAuthorErr             = "Unable to initialize TxAuthor"
	MaxAmntErr                  = "Unable to get max amount in account"
	NotConnectedErr             = "Not Connected To Decred Network"
	ParseFloatErr               = "Could not parse float"
	PasswordPopupIconsErr       = "Unable to load password popup icons"
	QrEncodeErr                 = "Could not generate encoded QR Image"
	ReceivePageLoadErr          = "Unable to load Receive Page components"
	SelectedWalletInvalidErr    = "Selected self sending wallet is invalid"
	SendPageLoadErr             = "Unable to load Send Page components"
	TransactionFeeSizeErr       = "Could not retrieve transaction fee and size"
	TransactionDetailsIconErr   = "Could not load transaction details icons"
	WrongPasswordErr            = "Wrong spending password. Please try again"
	WalletsErr                  = "Could not retrieve wallets"
	LoadMultiWalletErr          = "Unable to load MultiWallet"
)
