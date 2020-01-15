package wallet

import (
	"gioui.org/layout"
	//"gioui.org/unit"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets/security"
)

type (
	CreateWalletPage struct {
		multiWallet       *helper.MultiWallet
		changePageFunc    func(string)

		pinAndPasswordWidget *security.PinAndPasswordWidget

		seedPage          *SeedPage
		seed              string
		isShowingSeedPage bool
		isCreating        bool
		err               error
		validatePasswordErr error
	}
)

func NewCreateWalletPage(multiWallet *helper.MultiWallet) *CreateWalletPage {
	c := &CreateWalletPage{
		multiWallet:       multiWallet,
		isShowingSeedPage: false,
	}

	c.pinAndPasswordWidget = security.NewPinAndPasswordWidget(c.cancel, c.create)
	c.seedPage = NewSeedPage(c)

	return c
}

func (w *CreateWalletPage) GetWidgets(ctx *layout.Context, changePageFunc func(page string)) []func() {
	w.changePageFunc = changePageFunc 

	if w.isShowingSeedPage {
		w.seedPage.prepare(w.seed)
		return w.seedPage.getWidgets(ctx, changePageFunc)
	} 

	widgets := []func(){
		func(){
			helper.Inset(ctx, 20, helper.StandaloneScreenPadding, 0, helper.StandaloneScreenPadding, func(){
				w.pinAndPasswordWidget.Render(ctx)
			})
		},
	}

	return widgets
}

func (w *CreateWalletPage) cancel() {
	w.pinAndPasswordWidget.Reset()
	w.changePageFunc("welcome")
}

func (w *CreateWalletPage) create() {
	w.pinAndPasswordWidget.IsCreating = true

	/**doneChan := make(chan bool)
	go func() {
		defer func() {
			doneChan <- true
		}()
		/**wallet, err := w.multiWallet.CreateNewWallet("public", w.pinAndPasswordWidget.Value(), 0)
		if err != nil {
			w.err = err
			return
		}
		w.seed = wallet.Seed
		w.multiWallet.RegisterWalletID(wallet.ID)
	}()

	<-doneChan**/
	w.pinAndPasswordWidget.IsCreating = false
	w.isShowingSeedPage = true
}

