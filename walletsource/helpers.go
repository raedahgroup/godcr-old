package walletsource

type TransactionDirection int8

const (
	// TransactionDirectionSent for transactions sent to external address(es) from wallet
	TransactionDirectionSent TransactionDirection = iota

	// TransactionDirectionReceived for transactions received from external address(es) into wallet
	TransactionDirectionReceived

	// TransactionDirectionTransferred for transactions sent from wallet to internal address(es)
	TransactionDirectionTransferred
)

