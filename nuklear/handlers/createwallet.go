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
	seed                 string
	hasStoredSeed        bool
	validationErrors     map[string]string
}

func (handler *CreateWalletHandler) BeforeRender() {
	handler.err = nil
	handler.isRendering = false

	handler.seed = ""

	handler.passwordInput.Flags = nucular.EditField
	handler.passwordInput.PasswordChar = '*'

	handler.confirmPasswordInput.Flags = nucular.EditField
	handler.confirmPasswordInput.PasswordChar = '*'
}

func (handler *CreateWalletHandler) Render(window *nucular.Window, wallet app.WalletMiddleware, changePage func(string)) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.seed, handler.err = wallet.GenerateNewWalletSeed()
	}

	if contentWindow := helpers.NewWindow("Create Wallet Window", window, nucular.WindowBorder); contentWindow != nil {
		helpers.SetFont(window, helpers.PageContentFont)

		if handler.err != nil {
			contentWindow.Row(30).Dynamic(1)
			contentWindow.LabelColored(handler.err.Error(), "LC", helpers.DangerColor)
		}

		contentWindow.Row(30).Dynamic(1)
		contentWindow.Label("Create Wallet", "LC")

		contentWindow.Row(10).Dynamic(2)
		contentWindow.Label("Wallet Password", "LC")
		contentWindow.Label("Confirm Password", "LC")

		contentWindow.Row(20).Dynamic(2)
		handler.passwordInput.Edit(contentWindow.Window)
		handler.confirmPasswordInput.Edit(contentWindow.Window)

		passwordError, ok := handler.validationErrors["password"]
		if ok {
			contentWindow.LabelColored(passwordError, "LC", helpers.DangerColor)
		}

		if confirmPasswordError, ok := handler.validationErrors["confirmpassword"]; ok {
			if passwordError != "" {
				contentWindow.Label("", "LC")
			}
			contentWindow.LabelColored(confirmPasswordError, "LC", helpers.DangerColor)
		}

		contentWindow.Row(20).Dynamic(1)
		contentWindow.Label("Wallet Seed", "LC")

		contentWindow.Row(110).Dynamic(1)
		if seedWindow := helpers.NewWindow("Seed Window", contentWindow.Window, nucular.WindowBorder); seedWindow != nil {
			seedWindow.Row(80).Dynamic(1)
			seedWindow.LabelWrap(handler.seed)
			seedWindow.End()
		}

		contentWindow.Row(50).Dynamic(1)
		contentWindow.LabelWrapColored(`IMPORTANT: Keep the seed in a safe place as you will NOT be able to restore your wallet without it. Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.`,
			helpers.DangerColor,
		)

		contentWindow.Row(30).Dynamic(2)
		contentWindow.CheckboxText("I've stored the seed in a safe and secure location", &handler.hasStoredSeed)
		if hasStoredSeedError, ok := handler.validationErrors["hasstoredseed"]; ok {
			contentWindow.LabelColored("("+hasStoredSeedError+")", "LC", helpers.DangerColor)
		}

		contentWindow.Row(30).Static(200)
		if contentWindow.Button(label.T("Create Wallet"), false) {
			if !handler.hasErrors() {
				handler.err = wallet.CreateWallet(string(handler.passwordInput.Buffer), handler.seed)
				changePage("sync")
			}
			contentWindow.Master().Changed()
		}

		contentWindow.End()
	}
}

func (handler *CreateWalletHandler) hasErrors() bool {
	handler.validationErrors = make(map[string]string)

	password := string(handler.passwordInput.Buffer)
	confirmPassword := string(handler.confirmPasswordInput.Buffer)
	hasStoredSeed := handler.hasStoredSeed

	if password == "" {
		handler.validationErrors["password"] = "Wallet password is required"
	}

	if password != "" && confirmPassword != "" && password != confirmPassword {
		handler.validationErrors["confirmpassword"] = "Both passwords do not match"
	}

	if !hasStoredSeed {
		handler.validationErrors["hasstoredseed"] = "Please store seed and check this box"
	}

	return len(handler.validationErrors) > 0
}
