package helpers

import (
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/assets"
)

func TxDirectionIcon(direction int32) string {
	switch direction {
	case 0:
		return assets.SendIcon
	case 1:
		return assets.ReceiveIcon
	case 2:
		return assets.ReceiveIcon
	default:
		return assets.InfoIcon
	}
}

func ErrorHandler(err string, errorLabel *widget.Label) {
	errorLabel.SetText(err)
	errorLabel.Show()
}
