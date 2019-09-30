package pagehandlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

const (
	estimatedSettingsWindowHeight = 140
	estimatedGroupWindowPadding   = 12
)

type AccountsHandler struct {
	wallet             *dcrlibwallet.LibWallet
	err                error
	accounts           []*dcrlibwallet.Account
	isFetchingAccounts bool
	networkHDPath      string
	tickIcon           string
	crossIcon          string
}

func (handler *AccountsHandler) BeforeRender(wallet *dcrlibwallet.LibWallet, refreshWindowDisplay func()) {
	handler.wallet = wallet
	handler.err = nil
	handler.accounts = nil
	handler.isFetchingAccounts = false

	if handler.tickIcon == "" {
		handler.tickIcon = getTickIcon()
	}

	if handler.crossIcon == "" {
		handler.crossIcon = getCrossIcon()
	}

	if handler.networkHDPath == "" {
		if wallet.NetType() == "testnet3" {
			handler.networkHDPath = dcrlibwallet.TestnetHDPath
		} else {
			handler.networkHDPath = dcrlibwallet.MainnetHDPath
		}
	}
}

func getTickIcon() string {
	tickUnicode := "\\U2713"
	tickUnicodeInt, _ := strconv.ParseInt(strings.TrimPrefix(tickUnicode, "\\U"), 16, 32)

	return fmt.Sprintf("%s\n", string(tickUnicodeInt))
}

func getCrossIcon() string {
	crossUnicode := "\\U03c7"
	crossUnicodeInt, _ := strconv.ParseInt(strings.TrimPrefix(crossUnicode, "\\U"), 16, 32)

	return fmt.Sprintf("%s\n", string(crossUnicodeInt))
}

func (handler *AccountsHandler) Render(window *nucular.Window) {
	if handler.accounts == nil && handler.err == nil && !handler.isFetchingAccounts {
		go handler.fetchAccounts(window.Master().Changed)
	}

	widgets.PageContentWindowDefaultPadding("Accounts", window, func(contentWindow *widgets.Window) {
		// show accounts if any
		if len(handler.accounts) > 0 {
			handler.renderAccounts(contentWindow)
		}

		// show error if any
		if handler.err != nil {
			contentWindow.DisplayErrorMessage("", handler.err)
		}

		// show loading indicator if fetching
		if handler.isFetchingAccounts {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
}

func (handler *AccountsHandler) fetchAccounts(refreshWindowDisplay func()) {
	handler.isFetchingAccounts = true
	refreshWindowDisplay()

	defer func() {
		handler.isFetchingAccounts = false
		refreshWindowDisplay()
	}()

	getAccountsResp, err := handler.wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		handler.err = err
	} else {
		handler.accounts = getAccountsResp.Acc
	}
}

func (handler *AccountsHandler) renderAccounts(window *widgets.Window) {
	for index, account := range handler.accounts {
		headerLabel := account.Name + ":   " + dcrutil.Amount(account.Balance.Total).String()
		if account.Balance.Total != account.Balance.Spendable {
			headerLabel += fmt.Sprintf(" (Spendable: %s )", dcrutil.Amount(account.Balance.Spendable).String())
		}

		table := widgets.NewTable()
		table.AddRow(
			widgets.NewLabelTableCell("Account Number:", "LC"),
			widgets.NewLabelTableCell(strconv.Itoa(int(account.Number)), "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("HD Path:", "LC"),
			widgets.NewLabelTableCell(fmt.Sprintf("%s %d", handler.networkHDPath, account.Number), "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Keys:", "LC"),
			widgets.NewLabelTableCell(fmt.Sprintf("%d External, %d Internal, %d Imported",
				account.ExternalKeyCount, account.InternalKeyCount, account.ImportedKeyCount), "LC"),
		)

		tableHeight := table.Height()
		totalWindowPadding := estimatedGroupWindowPadding * 3

		window.Row(30).Dynamic(1)
		if window.TreePush(nucular.TreeNode, headerLabel, false) {
			window.Row(tableHeight + estimatedSettingsWindowHeight + totalWindowPadding).Dynamic(1)
			widgets.NoScrollGroupWindow(fmt.Sprintf("properties-window-%d", account.Number), window.Window, func(mainWindow *widgets.Window) {
				mainWindow.AddLabelWithFont("Properties", "LC", styles.BoldPageContentFont)

				mainWindow.Row(tableHeight + estimatedGroupWindowPadding).Dynamic(1)
				widgets.GroupWindow(fmt.Sprintf("table-window-%d", account.Number), mainWindow.Window, 0, func(tableWindow *widgets.Window) {
					table.Render(tableWindow)
				})
			})
			window.TreePop()
		}

		if index != len(handler.accounts)-1 {
			window.AddHorizontalLine(1, styles.BorderColor)
		}
	}
}
