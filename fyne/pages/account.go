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

var accountPageContainer pageContainer

func accountPage(wallet godcrApp.WalletMiddleware, appSettings *config.Settings, window fyne.Window) {
	titleLabel := widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	successLabel := widget.NewLabel("")
	errorLabel := widget.NewLabel("")
	errorLabel.Hide()

	accountNames, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		accountPageContainer.container.Children = []fyne.CanvasObject{widget.NewLabel("error getting accounts: " + err.Error())}
		widget.Refresh(accountPageContainer.container)
		return
	}

	listAccounts := receiveAccountDetails(accountNames, appSettings, wallet, errorLabel, successLabel)

	addAccountIcon := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		popup := createAccount(wallet, appSettings, listAccounts, window)
		popup.Show()
	})
	addAccount := widget.NewToolbar(addAccountIcon)
	container := widget.NewScrollContainer(listAccounts)

	container.Resize(fyne.NewSize(container.MinSize().Width, 500))

	output := widget.NewVBox(
		widget.NewHBox(titleLabel, addAccount),
		successLabel,
		errorLabel,
		fyne.NewContainer(container))

	accountPageContainer.container.Children = widget.NewHBox(widgets.NewHSpacer(10), output).Children
	widget.Refresh(accountPageContainer.container)
}

func createAccount(wallet godcrApp.WalletMiddleware, appSettings *config.Settings, listAccounts *widget.ScrollContainer, window fyne.Window) fyne.CanvasObject {
	// popUp houses all widgets, to display on account creation
	var popUp *widget.PopUp
	var createAccountButton *widget.Button

	titleLabel := widget.NewLabelWithStyle("Create new account", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	errorLabel := widget.NewLabel("")
	errorLabel.Hide()
	errorLabel.Hide()
	info := canvas.NewText("Accounts CANNOT be deleted once created", color.RGBA{255, 0, 0, menu.alphaTheme})
	info.TextStyle = fyne.TextStyle{Bold: true}

	name := widget.NewEntry()
	password := widget.NewPasswordEntry()
	name.SetPlaceHolder("Account name")
	password.SetPlaceHolder("Password")

	password.OnChanged = func(s string) {
		// disable button till there's a name and password.
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
			errorLabel.Show()
			errorLabel.SetText("could not create a new account: " + err.Error())
			return
		}

		// Reset the page.
		accountPage(wallet, appSettings, window)
		// Get the first object in the vbox container of account page.
		a, ok := interface{}(accountPageContainer.container.Children[1]).(*widget.Box)
		if !ok {
			return
		}

		successLabel, ok := interface{}(a.Children[1]).(*widget.Label)
		if !ok {
			return
		}
		successLabel.Show()
		successLabel.SetText("Account created")

		popUp.Hide()
	})
	createAccountButton.Disable()

	cancel := widget.NewButton("Cancel", func() {
		popUp.Hide()
	})

	output := widget.NewVBox(
		titleLabel,
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

func receiveAccountDetails(accounts []*walletcore.Account, appSettings *config.Settings, wallet godcrApp.WalletMiddleware, errorLabel, successLabel *widget.Label) *widget.ScrollContainer {
	scrollContainer := widget.NewScrollContainer(nil)
	collapsibleContainer := widget.NewVBox()
	propertiesContainer := make(map[string]*widget.Box)
	button := make(map[string]*widget.Button)
	container := make(map[string]*widget.Box)
	hideAccount := make(map[string]*widget.Check)
	defaultAccount := make(map[int]*widget.Check)
	accountNo := make(map[string]int)

	for i, acct := range accounts {
		accountName := acct.Name + ": " + acct.Balance.Total.String() + " (Spendable: " + acct.Balance.Spendable.String() + ")"
		accountNo[accountName] = i
		propertiesForm := widget.NewForm()

		propertiesForm.Append("Account Number", widget.NewLabelWithStyle(strconv.Itoa(int(acct.Number)), fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}))
		if wallet.NetType() == "testnet3" {
			propertiesForm.Append("HD Path", widget.NewLabelWithStyle(fmt.Sprintf("%s%d'", walletcore.TestnetHDPath, acct.Number), fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}))
		} else {
			propertiesForm.Append("HD Path", widget.NewLabelWithStyle(fmt.Sprintf("%s%d'", walletcore.MainnetHDPath, acct.Number), fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}))
		}
		keys := fmt.Sprintf("%d external, %d internal, %d imported", acct.ExternalKeyCount, acct.InternalKeyCount, acct.ImportedKeyCount)
		propertiesForm.Append("Keys", widget.NewLabelWithStyle(keys, fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}))

		hideAccount[accountName] = widget.NewCheck("Hide Account", func(s bool) {
			if s == true {
				// Filter for duplicate hidden account.
				hiddenAcct := make(map[uint32]uint32)
				for _, hidden := range appSettings.HiddenAccounts {
					hiddenAcct[hidden] = hidden
				}
				hiddenAcct[uint32(i)] = uint32(i)
				var hidden []uint32
				for value := range hiddenAcct {
					hidden = append(hidden, value)
				}
				appSettings.HiddenAccounts = hidden
				appSettings.HiddenAccounts = append(appSettings.HiddenAccounts, uint32(accountNo[accountName]))

			} else {
				var hiddenAccount []uint32
				for _, hidden := range appSettings.HiddenAccounts {
					if hidden == uint32(i) {
						continue
					}
					hiddenAccount = append(hiddenAccount, hidden)
				}
			}

			err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
				cnfg.HiddenAccounts = appSettings.HiddenAccounts
			})
			if err != nil {
				errorLabel.SetText("could not store hidden accounts to config " + err.Error())
				errorLabel.Show()
			}
		})

		for _, hidden := range appSettings.HiddenAccounts {
			if hidden == uint32(i) {
				fmt.Println("Hiden", appSettings.HiddenAccounts)
				hideAccount[accountName].SetChecked(true)
			}
		}

		(defaultAccount)[accountNo[accountName]] = widget.NewCheck("Default Account", func(s bool) {
			// Enable as default account and disable default account checkbox and hidden account checkbox
			// else enable hidden account checkbox.
			if s == true {
				// Remove account from hidden.
				if hideAccount[accountName].Checked {
					hideAccount[accountName].SetChecked(false)
					var hiddenAccounts []uint32
					for _, hidden := range appSettings.HiddenAccounts {
						if hidden == uint32(i) {
							continue
						}
						hiddenAccounts = append(hiddenAccounts, hidden)
					}
					appSettings.HiddenAccounts = hiddenAccounts
				}

				appSettings.DefaultAccount = uint32(accountNo[accountName])
				err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
					cnfg.HiddenAccounts = appSettings.HiddenAccounts
					cnfg.DefaultAccount = appSettings.DefaultAccount
				})
				if err != nil {
					errorLabel.SetText("could not store config " + err.Error())
					errorLabel.Show()
				}

				// Remove other account that has been set as default.
				for no := range accounts {
					if uint32(no) == appSettings.DefaultAccount {
						continue
					}
					defaultAccount[no].SetChecked(false)
				}
				defaultAccount[accountNo[accountName]].Disable()
				hideAccount[accountName].Disable()

			} else {
				hideAccount[accountName].Enable()
				defaultAccount[accountNo[accountName]].Enable()
			}
		})

		propertiesForm.Append("Hide Account", hideAccount[accountName])
		propertiesForm.Append("Default Account", defaultAccount[accountNo[accountName]])

		propertiesContainer[accountName] = widget.NewVBox(
			widget.NewLabelWithStyle("Wallet Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Properties", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			propertiesForm,
			widgets.NewVSpacer(10))
		propertiesContainer[accountName].Hide()

		button[accountName] = widget.NewButton(accountName, func() {
			if propertiesContainer[accountName].Hidden {
				propertiesContainer[accountName].Show()
				container[accountName].Children = []fyne.CanvasObject{button[accountName], propertiesContainer[accountName]}
			} else {
				propertiesContainer[accountName].Hide()
				container[accountName].Children = []fyne.CanvasObject{button[accountName]}
			}
			widget.Refresh(container[accountName])
			widget.Refresh(collapsibleContainer)
			scrollContainer.Resize(fyne.NewSize(500, 500))
			widget.Refresh(scrollContainer)
			widget.Refresh(accountPageContainer.container)
		})

		container[accountName] = widget.NewVBox()
		container[accountName].Append(button[accountName])
		collapsibleContainer.Append(container[accountName])
	}
	for i := range accounts {
		if appSettings.DefaultAccount == uint32(i) {
			defaultAccount[i].SetChecked(true)
			break
		}
	}
	scrollContainer.Content = widget.NewHBox(collapsibleContainer)
	return scrollContainer
}
