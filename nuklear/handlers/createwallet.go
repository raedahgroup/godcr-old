package handlers

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
	passwordInput        *nucular.TextEditor
	confirmPasswordInput *nucular.TextEditor
	seedBox              *nucular.TextEditor // todo why editor?
	seed                 string
	hasStoredSeed        bool
	validationErrors     map[string]string
}

func (handler *CreateWalletHandler) BeforeRender() {
	handler.err = nil
	handler.isRendering = false

	handler.seed = ""

	handler.passwordInput = &nucular.TextEditor{}
	handler.passwordInput.Flags = nucular.EditField
	handler.passwordInput.PasswordChar = '*'

	handler.confirmPasswordInput = &nucular.TextEditor{}
	handler.confirmPasswordInput.Flags = nucular.EditField
	handler.confirmPasswordInput.PasswordChar = '*'

	handler.validationErrors = make(map[string]string)
}

func (handler *CreateWalletHandler) Render(window *nucular.Window, wallet app.WalletMiddleware, changePage func(*nucular.Window, string)) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.seed, handler.err = wallet.GenerateNewWalletSeed()
	}

	widgets.PageContentWindowDefaultPadding("Create Wallet", window, func(contentWindow *widgets.Window) {
		columnWidths := []int{220, 220} // use 220 to hold each column to accommodate error messages
		contentWindow.AddLabelsWithWidths(columnWidths,
			widgets.NewLabelTableCell("Wallet Password", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Confirm Password", widgets.LeftCenterAlign),
		)

		contentWindow.AddEditorsWithWidths(columnWidths, handler.passwordInput, handler.confirmPasswordInput)

		passwordError, hasPasswordError := handler.validationErrors["password"]
		confirmPasswordError, hasConfirmPasswordError := handler.validationErrors["confirmpassword"]
		if hasPasswordError || hasConfirmPasswordError {
			formValidationLabels := make([]*widgets.LabelTableCell, 2)
			if hasPasswordError {
				formValidationLabels[0] = widgets.NewColoredLabelTableCell(passwordError, widgets.LeftCenterAlign,
					styles.DecredOrangeColor)
			}
			if hasConfirmPasswordError {
				formValidationLabels[1] = widgets.NewColoredLabelTableCell(confirmPasswordError, widgets.LeftCenterAlign,
					styles.DecredOrangeColor)
			}
			contentWindow.AddLabelsWithWidths(columnWidths, formValidationLabels...)
		}

		contentWindow.AddHorizontalSpace(20)
		contentWindow.AddLabelWithFont("Wallet Seed", widgets.LeftCenterAlign, styles.BoldPageContentFont)
		contentWindow.AddWrappedLabel(handler.seed, widgets.LeftCenterAlign)

		contentWindow.AddHorizontalSpace(10)
		contentWindow.AddWrappedLabelWithColor(walletcore.StoreSeedWarningText, widgets.LeftCenterAlign, styles.DecredOrangeColor)

		contentWindow.AddHorizontalSpace(10)
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

		contentWindow.AddHorizontalSpace(20)
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

		if handler.err != nil {
			contentWindow.DisplayErrorMessage("Error creating wallet", handler.err)
		}
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

	if password != "" && password != confirmPassword {
		handler.validationErrors["confirmpassword"] = "Both passwords do not match"
	}

	if !hasStoredSeed {
		handler.validationErrors["hasstoredseed"] = "Please store seed and check this box"
	}

	return len(handler.validationErrors) > 0
}
