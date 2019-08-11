package pages

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type accountPageData struct {
	errorLabel *canvas.Text
}

var account accountPageData

func accountPage(wallet godcrApp.WalletMiddleware, appSettings *config.Settings, window fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	successLabel := canvas.NewText("", color.RGBA{11, 156, 49, menu.alphaTheme})

	account.errorLabel = canvas.NewText("", color.RGBA{})
	account.errorLabel = canvas.NewText("", color.RGBA{255, 0, 0, menu.alphaTheme})
	account.errorLabel.Alignment = fyne.TextAlignCenter
	account.errorLabel.TextStyle = fyne.TextStyle{Bold: true}
	account.errorLabel.Hide()
	accountNames, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		//TODO: treat error here
	}

	listAccounts := receiveAccountDetails(accountNames, appSettings)

	addAccountIcon := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		popup := createAccount(wallet, appSettings, listAccounts, successLabel, window)
		popup.Show()
	})
	addAccount := widget.NewToolbar(addAccountIcon)

	container := widget.NewScrollContainer(listAccounts)
	fmt.Println("Container", container.MinSize())

	output := widget.NewVBox(
		widget.NewHBox(label, addAccount),
		successLabel,
		account.errorLabel,
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(container.MinSize().Width, container.MinSize().Height+300)), container))

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

//accountProperties creates a popUp that asks for account name and password so as to create the new account
func createAccount(wallet godcrApp.WalletMiddleware, appSettings *config.Settings, listAccounts *widget.Box, successLabel *canvas.Text, window fyne.Window) fyne.CanvasObject {
	//popUp houses all widgets, to display account creation
	var popUp *widget.PopUp
	var createAccountButton *widget.Button

	label := widget.NewLabelWithStyle("Create new account", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	info := canvas.NewText("Accounts CANNOT be deleted once created", color.RGBA{255, 0, 0, menu.alphaTheme})
	errorLabel := canvas.NewText("", color.RGBA{255, 0, 0, menu.alphaTheme})

	name := widget.NewEntry()
	password := widget.NewPasswordEntry()
	name.SetPlaceHolder("Account name")
	password.SetPlaceHolder("Password")

	password.OnChanged = func(s string) {
		if name.Text != "" && password.Text != "" {
			if createAccountButton.Disabled() {
				createAccountButton.Enable()
			}
		} else {
			if !createAccountButton.Disabled() {
				createAccountButton.Disable()
			}
		}
	}

	createAccountButton = widget.NewButton("Create", func() {
		_, err := wallet.NextAccount(name.Text, password.Text)
		if err != nil {
			errorLabel.Text = "could not create new account " + err.Error()
			canvas.Refresh(errorLabel)
		} else {
			//if there were recent
			successLabel.Text = "Account created"
			canvas.Refresh(successLabel)

			accountNames, _ := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
			listAccounts = receiveAccountDetails(accountNames, appSettings)
			widget.Refresh(listAccounts)
			popUp.Hide()
		}
	})

	cancel := widget.NewButton("Cancel", func() {
		name.SetText("")
		password.SetText("")
		popUp.Hide()
	})

	output := widget.NewVBox(
		label,
		widgets.NewVSpacer(10),
		info,
		errorLabel,
		name,
		password,
		widget.NewHBox(layout.NewSpacer(), createAccountButton, widgets.NewHSpacer(20), cancel, layout.NewSpacer()),
	)
	popUp = widget.NewModalPopUp(output, window.Canvas())
	return popUp
}

// func receiveAccountDetails(accounts []*walletcore.Account, appSettings config.Settings) *widget.Box {
// 	overallContainer := widget.NewVBox()
// 	var walletInfoContainer []fyne.CanvasObject

// 	propertiesForm := widget.NewForm()
// 	walletSettingsForm := widget.NewForm()
// 	propertiesForm.Append("Account Number", widget.NewLabelWithStyle(strconv.Itoa(int(account.Number)), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}))
// 	propertiesForm.Append("HD Path", widget.NewLabelWithStyle(strconv.Itoa(int(account.Number)), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}))
// 	walletSettingsForm.Append()

// 	//initially hide form

// 	var button *widget.Button
// 	var buttons []fyne.CanvasObject

// 	for i, account := range accounts {
// 		container := widget.NewVBox()

// 		container.Append(widget.NewLabelWithStyle("Properties", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
// 		container.Append(propertiesForm)
// 		container.Append(widgets.NewVSpacer(10))
// 		container.Append(widget.NewLabelWithStyle("Wallet Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
// 		container.Append(walletSettingsForm)
// 		container.Hide()

// 		append(buttons, widget.NewButton(account.Name, func() {
// 			//TODO: remove this when issue #377 is merged to masters
// 			if (*walletInfoContainer)[i].Hidden {
// 				overallContainer.Children = []fyne.CanvasObject{buttons[i], walletInfoContainer}
// 				(*walletInfoContainer)[i].Show()
// 			} else {
// 				overallContainer.Children = []fyne.CanvasObject{button}
// 				walletInfoContainer.Hide()
// 			}
// 			widget.Refresh(overallContainer)
// 		}))
// 	}

// 	overallContainer.Append(button)
// 	return overallContainer
// }

func receiveAccountDetails(accounts []*walletcore.Account, appSettings *config.Settings) *widget.Box {
	fmt.Println(appSettings)
	var defaultAccount []*widget.Check
	defaultAccount = make([]*widget.Check, len(accounts))
	container := widget.NewVBox()

	for i, account := range accounts {
		propertiesForm := widget.NewForm()
		propertiesForm.Append("Account Number", widget.NewLabelWithStyle(strconv.Itoa(int(account.Number)), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}))
		propertiesForm.Append("HD Path", widget.NewLabelWithStyle(strconv.Itoa(int(account.Number)), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}))

		fmt.Println(appSettings.HiddenAccounts)
		hideAccount := widget.NewCheck("Hide Account", func(s bool) {
			if s == true {

			}
			appSettings.HiddenAccounts = append(appSettings.HiddenAccounts, uint32(i))
		})

		for _, hidden := range appSettings.HiddenAccounts {
			if hidden == uint32(i) {
				fmt.Println("Hiden", appSettings.HiddenAccounts)
				hideAccount.SetChecked(true)
			}
		}

		(defaultAccount)[i] = widget.NewCheck("Default Account", func(s bool) {
			//remove all accounts that are defaults
			//check if account is hidden, if hidden remove from being hidden and disable hidden check
			if hideAccount.Checked && s == true {
				hideAccount.SetChecked(false)
				hideAccount.Disable()
				var hiddenAccountNo []uint32
				for _, hidden := range appSettings.HiddenAccounts {
					if hidden == uint32(i) {
						continue
					}
					hiddenAccountNo = append(hiddenAccountNo, hidden)
				}
				appSettings.HiddenAccounts = hiddenAccountNo
			} else {
				if hideAccount.Disabled() {
					hideAccount.Enable()
				}
			}

			appSettings.DefaultAccount = uint32(i)
		})
		if appSettings.DefaultAccount == uint32(i) {
			defaultAccount[i].SetChecked(true)
		}

		propertiesForm.Append("Hide Account", hideAccount)
		propertiesForm.Append("Default Account", defaultAccount[i])

		propertiesContainer := widget.NewVBox(
			widget.NewLabelWithStyle("Wallet Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Properties", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			propertiesForm,
			widgets.NewVSpacer(10))

		button := widget.NewButton(account.Name+": "+account.Balance.Total.String()+" (Spendable: "+account.Balance.Spendable.String()+")", func() {
			if propertiesContainer.Hidden {
				propertiesContainer.Show()
				//container.Children = []fyne.CanvasObject{fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(1000, 30)), button), propertiesContainer}

			} else {
				propertiesContainer.Hide()
				//container.Children = []fyne.CanvasObject{fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(1000, 30)), button)}
			}
			widget.Refresh(container)
		})

		container.Append(fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(1000, 30)), button))
		container.Append(propertiesContainer)
	}

	return container
}
