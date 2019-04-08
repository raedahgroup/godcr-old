package helpers

import (
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func RequestSpendingPassphrase(pages *tview.Pages, successFunc func(string), cancelFunc func()) {
	passphraseField := tview.NewInputField().SetMaskCharacter('*')
	passphraseModal := primitives.NewFormModal("Enter Spending Passphrase").
		AddFormItem(passphraseField)

	closeModal := func() {
		pages.RemovePage("passphrase")
		passphraseField.SetText("") // clear sensitive info
	}

	// wrapper around passed in cancelFunc to closse modal before invoking provided cancelFunc
	cancelFuncWrapper := func() {
		closeModal()
		cancelFunc()
	}

	passphraseModal.AddButton("Submit", func() {
		passphrase := passphraseField.GetText()
		// close modal before calling success func so that the modal does not remain open if success func takes some time to execute
		closeModal()
		successFunc(passphrase)
	})

	passphraseModal.AddButton("Cancel", cancelFuncWrapper)
	passphraseModal.SetCancelFunc(cancelFuncWrapper)

	pages.AddPage("passphrase", passphraseModal, true, true)
}