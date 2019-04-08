package helpers

import (
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func RequestSpendingPassphrase(pages *tview.Pages, successFunc func(string), cancelFunc func()) {
	passphraseField := tview.NewInputField().SetMaskCharacter('*')
	passphraseModal := primitives.NewFormModal("Enter Spending Passphrase").
		AddFormItem(passphraseField)

	passphraseModal.AddButton("Submit", func() {
		pages.RemovePage("passphrase")
		passphraseField.SetText("") // clear sensitive info
		successFunc(passphraseField.GetText())
	})
	passphraseModal.AddButton("Cancel", func() {
		pages.RemovePage("passphrase")
		passphraseField.SetText("") // clear sensitive info
		cancelFunc()
	})
	passphraseModal.SetCancelFunc(func() {
		pages.RemovePage("passphrase")
		passphraseField.SetText("") // clear sensitive info
		cancelFunc()
	})

	pages.AddPage("passphrase", passphraseModal, true, true)
}