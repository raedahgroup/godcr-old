package pagehandlers

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
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

func (handler *CreateWalletHandler) Render(window *nucular.Window, wallet app.WalletMiddleware, changePage func(*nucular.Window, string)) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.seed, handler.err = wallet.GenerateNewWalletSeed()
	}

	widgets.PageContentWindowDefaultPadding("Create Wallet", window, func(contentWindow *widgets.Window) {
		if handler.err != nil {
			contentWindow.DisplayErrorMessage("Error creating wallet", handler.err)
		}

		createWalletForm := widgets.NewTable()

		createWalletForm.AddRow(
			widgets.NewLabelTableCell("Wallet Password", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Confirm Password", widgets.LeftCenterAlign),
		)

		createWalletForm.AddRow(
			widgets.NewEditTableCell(handler.passwordInput, 200),
			widgets.NewEditTableCell(handler.confirmPasswordInput, 200),
		)

		formValidationCells := make([]widgets.TableCell, 2)
		if passwordError, hasPasswordError := handler.validationErrors["password"]; hasPasswordError {
			formValidationCells[0] = widgets.NewColoredLabelTableCell(passwordError, widgets.LeftCenterAlign, styles.DecredOrangeColor)
		}
		if confirmPasswordError, hasConfirmPasswordError := handler.validationErrors["confirmpassword"]; hasConfirmPasswordError {
			formValidationCells[1] = widgets.NewColoredLabelTableCell(confirmPasswordError, widgets.LeftCenterAlign, styles.DecredOrangeColor)
		}
		createWalletForm.AddRow(formValidationCells...)

		createWalletForm.Render(contentWindow)

		contentWindow.AddLabel("Wallet Seed", widgets.LeftCenterAlign)
		contentWindow.AddWrappedLabel(handler.seed, widgets.LeftCenterAlign)

		contentWindow.AddWrappedLabelWithColor(walletcore.StoreSeedWarningText, widgets.LeftCenterAlign, styles.DecredOrangeColor)

		contentWindow.AddCheckbox("I've stored the seed in a safe and secure location", &handler.hasStoredSeed, func() {
			if !handler.hasStoredSeed {
				handler.validationErrors["hasstoredseed"] = "Please store seed and check this box"
			} else {
				delete(handler.validationErrors, "hasstoredseed")
			}
			contentWindow.Master().Changed()
		})
		if hasStoredSeedError, ok := handler.validationErrors["hasstoredseed"]; ok {
			contentWindow.AddColoredLabel(hasStoredSeedError, styles.DecredOrangeColor, widgets.LeftCenterAlign)
		}

		contentWindow.AddButton("Create Wallet", func() {
			if !handler.hasErrors() {
				handler.err = wallet.CreateWallet(string(handler.passwordInput.Buffer), handler.seed)
				if handler.err != nil {
					changePage(window, "sync")
				} else {
					contentWindow.Master().Changed() // refresh to display error
				}
			}
		})
	})
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
