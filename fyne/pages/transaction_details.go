package pages

import (
	"context"
	"fmt"
	"strconv"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type TransactionDetails struct {
	err       error
	ctx       context.Context
	wallet    walletcore.Wallet
	hash      string
	container *widgets.Box
}

func NewTransactionDetailsHandler(hash string, ctx context.Context, wallet walletcore.Wallet) *TransactionDetails {
	return &TransactionDetails{
		err:    nil,
		ctx:    ctx,
		hash:   hash,
		wallet: wallet,
	}
}

func (t *TransactionDetails) Render(from string, container *widgets.Box) {
	t.container = container
	container.Empty()
	defer container.Update()

	// set main view title
	container.SetTitle("Transaction Details")

	// create subview
	view := widgets.NewVBox()
	t.addBreadcrumb(from, view)

	// fetch transaction details
	txDetails, err := t.wallet.GetTransaction(t.hash)
	if err != nil {
		container.AddLabel(err.Error())
		return
	}

	table := widgets.NewTable()
	table.AddRowSimple("Confirmations", strconv.Itoa(int(txDetails.Confirmations)))
	table.AddRowSimple("Hash", txDetails.Hash)
	table.AddRowSimple("Block Height", strconv.Itoa(int(txDetails.BlockHeight)))
	table.AddRowSimple("Direction", txDetails.Direction.String())
	table.AddRowSimple("Type", txDetails.Type)
	table.AddRowSimple("Amount", dcrutil.Amount(txDetails.Amount).String())
	table.AddRowSimple("Size", strconv.Itoa(txDetails.Size)+" Bytes")
	table.AddRowSimple("Fee", dcrutil.Amount(txDetails.Fee).String())
	table.AddRowSimple("Fee Rate", dcrutil.Amount(txDetails.FeeRate).String())
	table.AddRowSimple("Time", fmt.Sprintf("%s UTC", txDetails.LongTime))
	view.Add(table.CondensedTable())

	view.Add(widgets.NewVSpacer(20))
	view.AddBoldLabel("Inputs")

	inputsTable := widgets.NewTable()
	inputsTable.AddRowSimple("Previous Outpoint", "Amount")
	for _, input := range txDetails.Inputs {
		inputsTable.AddRowSimple(input.PreviousTransactionHash, dcrutil.Amount(input.Amount).String())
	}
	view.Add(inputsTable.CondensedTable())

	view.Add(widgets.NewVSpacer(20))
	view.AddBoldLabel("Outputs")

	outputsTable := widgets.NewTable()
	outputsTable.AddRowSimple("Address", "Account", "Value", "Type")
	for _, output := range txDetails.Outputs {
		outputsTable.AddRowSimple(
			output.Address,
			output.AccountName,
			dcrutil.Amount(output.Amount).String(),
			output.ScriptType,
		)
	}
	view.Add(outputsTable.CondensedTable())

	// add subview to main view
	container.Add(view)

}

func (t *TransactionDetails) addBreadcrumb(from string, view *widgets.Box) {
	// add breadcrumb
	breadcrumb := []*widgets.Breadcrumb{
		{
			Text: from,
			Action: func() {
				previousPage := &OverviewHandler{}
				t.container.Empty()
				t.container.Update()
				previousPage.Render(t.ctx, t.wallet, t.container)
			},
		},
		{
			Text:   "Transaction Details",
			Action: nil,
		},
	}

	view.AddBreadcrumb(breadcrumb)
}
