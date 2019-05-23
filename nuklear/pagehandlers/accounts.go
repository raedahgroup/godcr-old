package pagehandlers

import (
	"errors"
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/command"
	"github.com/aarzilli/nucular/rect"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
	"golang.org/x/mobile/event/mouse"
)

const (
	estimatedSettingsWindowHeight = 140
	estimatedGroupWindowPadding   = 12
	strokeHeight                  = 1
	rectXPadding                  = 15
	rectYPadding                  = 10
)

type account struct {
	isSetAsHidden         bool
	isSetAsDefaultAccount bool
	account               *walletcore.Account
}

type AccountsHandler struct {
	err                error
	accounts           []*account
	wallet             walletcore.Wallet
	settings           *config.Settings
	isFetchingAccounts bool
	networkHDPath      string
	tickIcon           string
	crossIcon          string
}

func (handler *AccountsHandler) BeforeRender(wallet walletcore.Wallet, settings *config.Settings, refreshWindowDisplay func()) bool {
	handler.wallet = wallet
	handler.err = nil
	handler.accounts = nil
	handler.isFetchingAccounts = false
	handler.settings = settings

	if handler.tickIcon == "" {
		handler.tickIcon = getTickIcon()
	}

	if handler.crossIcon == "" {
		handler.crossIcon = getCrossIcon()
	}

	if handler.networkHDPath == "" {
		if wallet.NetType() == "testnet3" {
			handler.networkHDPath = walletcore.TestnetHDPath
		} else {
			handler.networkHDPath = walletcore.MainnetHDPath
		}
	}

	return true
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

	accounts, err := handler.wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		handler.err = err
		return
	}

	for _, accountItem := range accounts {
		var isSetAsHidden, isSetAsDefaultAccount bool
		for _, hiddenAccount := range handler.settings.HiddenAccounts {
			if uint32(hiddenAccount) == accountItem.Number {
				isSetAsHidden = true
				break
			}
		}

		if handler.settings.DefaultAccount == accountItem.Number {
			isSetAsDefaultAccount = true
		}

		acc := &account{
			isSetAsHidden:         isSetAsHidden,
			isSetAsDefaultAccount: isSetAsDefaultAccount,
			account:               accountItem,
		}
		handler.accounts = append(handler.accounts, acc)
	}
}

func (handler *AccountsHandler) renderAccounts(window *widgets.Window) {
	for index, item := range handler.accounts {
		headerLabel := item.account.Name + ":   " + item.account.Balance.Total.String()
		if item.account.Balance.Total != item.account.Balance.Spendable {
			headerLabel += fmt.Sprintf(" (Spendable: %s )", item.account.Balance.Spendable.String())
		}

		table := widgets.NewTable()
		table.AddRow(
			widgets.NewLabelTableCell("Account Number:", "LC"),
			widgets.NewLabelTableCell(strconv.Itoa(int(item.account.Number)), "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("HD Path:", "LC"),
			widgets.NewLabelTableCell(fmt.Sprintf("%s %d", handler.networkHDPath, item.account.Number), "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Keys:", "LC"),
			widgets.NewLabelTableCell(fmt.Sprintf("%d External, %d Internal, %d Imported", item.account.ExternalKeyCount, item.account.InternalKeyCount, item.account.ImportedKeyCount), "LC"),
		)

		tableHeight := table.Height()
		totalWindowPadding := estimatedGroupWindowPadding * 3

		window.Row(30).Dynamic(1)
		if window.TreePush(nucular.TreeNode, headerLabel, false) {
			window.Row(tableHeight + estimatedSettingsWindowHeight + totalWindowPadding).Dynamic(1)
			widgets.NoScrollGroupWindow(fmt.Sprintf("properties-window-%d", item.account.Number), window.Window, func(mainWindow *widgets.Window) {
				mainWindow.AddLabelWithFont("Properties", "LC", styles.BoldPageContentFont)

				mainWindow.Row(tableHeight + estimatedGroupWindowPadding).Dynamic(1)
				widgets.GroupWindow(fmt.Sprintf("table-window-%d", item.account.Number), mainWindow.Window, 0, func(tableWindow *widgets.Window) {
					table.Render(tableWindow)
				})

				mainWindow.AddLabelWithFont("Wallet Settings", "LC", styles.BoldPageContentFont)
				mainWindow.Row(estimatedSettingsWindowHeight + estimatedGroupWindowPadding).Dynamic(1)
				widgets.GroupWindow(fmt.Sprintf("settings-window-%d", item.account.Number), mainWindow.Window, 0, func(settingsWindow *widgets.Window) {
					settingsWindow.Row(estimatedSettingsWindowHeight - 30).Dynamic(1)
					bounds, out := settingsWindow.Custom(nstyle.WidgetStateInactive)
					handler.drawCustomCheckbox(item, settingsWindow, bounds, out)

				})
			})
			window.TreePop()
		}

		if index != len(handler.accounts)-1 {
			window.AddHorizontalLine(1, styles.BorderColor)
		}
	}
}

func (handler *AccountsHandler) drawCustomCheckbox(account *account, window *widgets.Window, bounds rect.Rect, out *command.Buffer) {
	accountVisibiltyRect, defaultAccountRect := drawRectangle(window, bounds, out)

	// account visibilty section
	accountVisibiltyRectInnerRect := rect.Rect{
		X: accountVisibiltyRect.X + strokeHeight,
		Y: accountVisibiltyRect.Y + strokeHeight,
		W: accountVisibiltyRect.W - (2 * strokeHeight),
		H: accountVisibiltyRect.H - (2 * strokeHeight),
	}

	accountVisibilityIcon := handler.tickIcon

	// listen account visibilty events
	if window.Input().Mouse.HoveringRect(accountVisibiltyRect) {
		out.FillRect(accountVisibiltyRectInnerRect, 0, styles.AlternateGrayColor)
	}

	if !account.isSetAsHidden {
		accountVisibilityIcon = handler.crossIcon
	}
	drawText(accountVisibiltyRect, out, "Hide account", "Account balance will be ignored", accountVisibilityIcon)

	if window.Input().Mouse.IsClickDownInRect(mouse.ButtonLeft, accountVisibiltyRect, false) {
		handler.toggleAccountVisibilty(account, window)
	}

	// default account section
	defaultAccountRectInnerRect := rect.Rect{
		X: defaultAccountRect.X + strokeHeight,
		Y: defaultAccountRect.Y,
		W: defaultAccountRect.W - (2 * strokeHeight),
		H: defaultAccountRect.H - (2 * strokeHeight),
	}

	defaultAccountIcon := handler.tickIcon
	if !account.isSetAsDefaultAccount {
		defaultAccountIcon = handler.crossIcon
	}

	if window.Input().Mouse.HoveringRect(defaultAccountRect) {
		out.FillRect(defaultAccountRectInnerRect, 0, styles.AlternateGrayColor)
	}

	drawText(defaultAccountRect, out, "Default account", "Make this account default for all outgoing and incoming transactions", defaultAccountIcon)

	if window.Input().Mouse.IsClickDownInRect(mouse.ButtonLeft, defaultAccountRect, false) {
		handler.toggleDefaultAccount(account, window)
	}
}

func (handler *AccountsHandler) toggleAccountVisibilty(accountItem *account, window *widgets.Window) {
	defer window.Master().Changed()

	var toggleAccountVisibilityFunc func(*account, *widgets.Window) error
	if accountItem.isSetAsHidden {
		toggleAccountVisibilityFunc = handler.revealAccount
	} else {
		toggleAccountVisibilityFunc = handler.hideAccount
	}

	if err := toggleAccountVisibilityFunc(accountItem, window); err != nil {
		widgets.NewAlertWidget(fmt.Sprintf("Error saving changes: %s", err.Error()), true, window)
		return
	}

	accountItem.isSetAsHidden = !accountItem.isSetAsHidden
	widgets.NewAlertWidget("Changes saved successfully!", false, window)
}

func (handler *AccountsHandler) hideAccount(accountItem *account, window *widgets.Window) error {
	accountToBeHidden := accountItem.account.Number
	hiddenAccounts := handler.settings.HiddenAccounts

	// make sure the account is not already set to be hidden
	for _, hiddenAccount := range hiddenAccounts {
		if hiddenAccount == accountToBeHidden {
			return errors.New("Account is already hidden")
		}
	}

	hiddenAccounts = append(hiddenAccounts, accountToBeHidden)
	err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
		cnfg.HiddenAccounts = hiddenAccounts
	})
	if err != nil {
		return err
	}
	handler.settings.HiddenAccounts = hiddenAccounts
	return nil
}

func (handler *AccountsHandler) revealAccount(accountItem *account, window *widgets.Window) error {
	hiddenAccounts := handler.settings.HiddenAccounts
	for index := range handler.settings.HiddenAccounts {
		if hiddenAccounts[index] == accountItem.account.Number {
			hiddenAccounts = append(hiddenAccounts[:index], hiddenAccounts[index+1:]...)
			err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
				cnfg.HiddenAccounts = hiddenAccounts
			})
			if err != nil {
				return err
			}

			handler.settings.HiddenAccounts = hiddenAccounts
			return nil
		}
	}
	return errors.New("Error revealing account")
}

func (handler *AccountsHandler) toggleDefaultAccount(accountItem *account, window *widgets.Window) {
	defer window.Master().Changed()

	fmt.Println("dddd")

	var toggleAccountFunc func(accountItem *account, window *widgets.Window) error
	if accountItem.isSetAsDefaultAccount {
		toggleAccountFunc = handler.unsetDefaultAccount
	} else {
		toggleAccountFunc = handler.setDefaultAccount
	}

	if err := toggleAccountFunc(accountItem, window); err != nil {
		widgets.NewAlertWidget(fmt.Sprintf("Error saving changes: %s", err.Error()), true, window)
		return
	}

	accountItem.isSetAsDefaultAccount = !accountItem.isSetAsDefaultAccount
	widgets.NewAlertWidget("Changes saved successfully!", false, window)
}

func (handler *AccountsHandler) setDefaultAccount(accountItem *account, window *widgets.Window) error {
	// first check if account is alreadr default
	if handler.settings.DefaultAccount == accountItem.account.Number {
		return errors.New("The account is already set as default")
	}

	// uncheck the other default account checkbox
	for _, acc := range handler.accounts {
		if acc.isSetAsDefaultAccount && acc.account.Number != accountItem.account.Number {
			acc.isSetAsDefaultAccount = false
		}
	}

	err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
		cnfg.DefaultAccount = accountItem.account.Number
	})
	if err != nil {
		return err
	}
	handler.settings.DefaultAccount = accountItem.account.Number
	return nil
}

func (handler *AccountsHandler) unsetDefaultAccount(accountItem *account, window *widgets.Window) error {
	// first check if account is default
	if handler.settings.DefaultAccount != accountItem.account.Number {
		return errors.New("This account is not set as default")
	}

	err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
		cnfg.DefaultAccount = 0
	})
	if err != nil {
		return err
	}
	handler.settings.DefaultAccount = 0
	return nil
}

func drawRectangle(window *widgets.Window, bounds rect.Rect, out *command.Buffer) (rect.Rect, rect.Rect) {
	topLeftPoint := bounds.Min()
	topRightPoint := image.Point{bounds.X + bounds.W - strokeHeight, bounds.Y}
	bottomLeftPoint := image.Point{bounds.X, bounds.Y + bounds.H - strokeHeight}
	bottomRightPoint := image.Point{bounds.X + bounds.W - strokeHeight, bounds.Y + bounds.H - strokeHeight}
	middleLeftPoint := image.Point{bounds.X, bounds.Y + (bounds.H / 2) - strokeHeight}
	middleRightPoint := image.Point{bounds.X + bounds.W - strokeHeight, bounds.Y + (bounds.H / 2) - strokeHeight}

	out.StrokeLine(topLeftPoint, topRightPoint, strokeHeight, styles.BorderColor)
	out.StrokeLine(topLeftPoint, bottomLeftPoint, strokeHeight, styles.BorderColor)
	out.StrokeLine(bottomLeftPoint, bottomRightPoint, strokeHeight, styles.BorderColor)
	out.StrokeLine(topRightPoint, bottomRightPoint, strokeHeight, styles.BorderColor)
	out.StrokeLine(middleLeftPoint, middleRightPoint, strokeHeight, styles.BorderColor)

	topRect := rect.Rect{
		X: topLeftPoint.X,
		Y: topLeftPoint.Y,
		W: bounds.W,
		H: bounds.H / 2,
	}

	bottomRect := rect.Rect{
		X: topLeftPoint.X,
		Y: topLeftPoint.Y + (bounds.H / 2),
		W: bounds.W,
		H: bounds.H / 2,
	}

	return topRect, bottomRect
}

func drawText(bounds rect.Rect, out *command.Buffer, mainText, leadText, icon string) {
	iconRect := rect.Rect{
		X: bounds.X + rectXPadding,
		Y: bounds.Y + rectYPadding - 3,
		W: 15,
		H: 20,
	}

	mainTextRect := rect.Rect{
		X: iconRect.X + 15,
		Y: bounds.Y + rectYPadding,
		W: bounds.W,
		H: 20,
	}

	leadTextRect := rect.Rect{
		X: bounds.X + rectXPadding,
		Y: mainTextRect.Y + mainTextRect.H,
		W: bounds.W,
		H: 20,
	}

	out.DrawText(iconRect, icon, styles.BoldPageContentFont, styles.BlackColor)
	out.DrawText(mainTextRect, mainText, styles.SmallBoldPageContentFont, styles.BlackColor)
	out.DrawText(leadTextRect, leadText, styles.PageContentFont, styles.GrayColor)

}
