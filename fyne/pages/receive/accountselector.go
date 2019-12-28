package receive

import (
	"time"

	"fyne.io/fyne"

	"github.com/raedahgroup/godcr/fyne/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (receivePage *ReceivePageObjects) initAccountSelector() error {
	receivePage.Accounts.OnAccountChange = receivePage.onAccountChange
	accountBox, err := receivePage.Accounts.CreateAccountSelector(values.ReceivingAccountLabel)
	if err != nil {
		return err
	}

	receivePage.ReceivePageContents.Append(accountBox)
	return err
}

func (receivePage *ReceivePageObjects) onAccountChange() {
	receivePage.generateAddressAndQR(false)
}

func (receivePage *ReceivePageObjects) showInfoLabel(object *widgets.BorderedText, err string) {
	object.SetText(err)
	object.SetPadding(fyne.NewSize(20, 8))
	object.Container.Show()
	receivePage.ReceivePageContents.Refresh()

	time.AfterFunc(time.Second*5, func() {
		object.Container.Hide()
		receivePage.ReceivePageContents.Refresh()
	})
}
