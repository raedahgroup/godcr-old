package pagehandlers

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type OverviewHandler struct {
	err         error
	isRendering bool
	accounts    []*walletcore.Account
	detailed    bool
}

func (handler *OverviewHandler) BeforeRender() {
	handler.err = nil
	handler.accounts = nil
	handler.isRendering = false
	handler.detailed = false
}

func (handler *OverviewHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.accounts, handler.err = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	}

	widgets.PageContentWindowDefaultPadding("Overview", window, func(contentWindow *widgets.Window) {
		contentWindow.AddLabelWithFont("Current Total Balance", widgets.LeftCenterAlign, styles.BoldPageContentFont)

		if handler.err != nil {
			contentWindow.DisplayErrorMessage(handler.err.Error())
		} else {
			contentWindow.AddLabel(walletcore.WalletBalance(handler.accounts), widgets.LeftCenterAlign)
		}
	})
}
