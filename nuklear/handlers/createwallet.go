package handlers

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type CreateWalletHandler struct {
	err                  error
	isRendering          bool
	passwordInput        nucular.TextEditor
	confirmPasswordInput nucular.TextEditor
	seedBox              nucular.TextEditor
	walletSeed           string
	hasStoredSeed        bool
}

func (handler *CreateWalletHandler) BeforeRender() {
	handler.err = nil
	handler.isRendering = false

	handler.walletSeed = ""

	handler.passwordInput.Flags = nucular.EditField
	handler.passwordInput.PasswordChar = '*'

	handler.confirmPasswordInput.Flags = nucular.EditField
	handler.confirmPasswordInput.PasswordChar = '*'
}

func (handler *CreateWalletHandler) Render(window *nucular.Window, wallet app.WalletMiddleware) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.walletSeed, handler.err = wallet.GenerateNewWalletSeed()
	}

	if contentWindow := helpers.NewWindow("Create Wallet Window", window, nucular.WindowBorder); contentWindow != nil {
		helpers.SetFont(window, helpers.PageContentFont)

		if handler.err != nil {
			contentWindow.SetErrorMessage(handler.err.Error())
		} else {
			contentWindow.Row(30).Dynamic(1)
			contentWindow.Label("Create Wallet", "LC")

			contentWindow.Row(10).Dynamic(2)
			contentWindow.Label("Wallet Password", "LC")
			contentWindow.Label("Confirm Password", "LC")

			contentWindow.Row(30).Dynamic(2)
			handler.passwordInput.Edit(contentWindow.Window)
			handler.confirmPasswordInput.Edit(contentWindow.Window)

			contentWindow.Row(20).Dynamic(1)
			contentWindow.Label("Wallet Seed", "LC")

			contentWindow.Row(100).Dynamic(1)
			if seedWindow := helpers.NewWindow("Seed Window", contentWindow.Window, nucular.WindowBorder); seedWindow != nil {
				seedWindow.Row(80).Dynamic(1)
				seedWindow.LabelWrap(handler.walletSeed)
				seedWindow.End()
			}

			contentWindow.Row(50).Dynamic(1)
			contentWindow.LabelWrapColored(`IMPORTANT: Keep the seed in a safe place as you will NOT be able to restore your wallet without it. Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.`,
				helpers.DangerColor,
			)

			contentWindow.Row(39).Dynamic(1)
			contentWindow.CheckboxText("I've stored the seed in a safe and secure location", &handler.hasStoredSeed)

			contentWindow.Row(30).Static(200)
			if contentWindow.Button(label.T("Create Wallet"), false) {

			}
		}

		contentWindow.End()
	}
}
