package handlers

import (
	"image"
	"image/draw"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/helpers"
	qrcode "github.com/skip2/go-qrcode"
)

type ReceiveHandler struct {
	err         error
	isRendering bool
	accounts    []*walletcore.Account

	// form selector index
	selectedAccountIndex  int
	selectedAccountNumber uint32

	// generatedAddress
	generatedAddress string
}

func (handler *ReceiveHandler) BeforeRender() {
	handler.err = nil
	handler.accounts = nil
	handler.isRendering = false

	// form selector index
	handler.selectedAccountIndex = 0
	handler.selectedAccountNumber = uint32(0)
}

func (handler *ReceiveHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.accounts, handler.err = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	}

	// draw page
	if pageWindow := helpers.NewWindow("Receive Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("Generate Receive Address")

		// content window
		if contentWindow := pageWindow.ContentWindow("Receive Content"); contentWindow != nil {
			if handler.err != nil {
				contentWindow.SetErrorMessage(handler.err.Error())
			} else {
				accountNames := make([]string, len(handler.accounts))
				for index, account := range handler.accounts {
					accountNames[index] = account.Name
				}

				contentWindow.Row(25).Static(300, 100)
				// draw select account combo
				handler.selectedAccountIndex = contentWindow.ComboSimple(accountNames, handler.selectedAccountIndex, 25)

				// draw submit button
				if contentWindow.Button(label.T("Generate"), false) {
					// get selected account by index
					accountName := accountNames[handler.selectedAccountIndex]
					for _, account := range handler.accounts {
						if account.Name == accountName {
							handler.selectedAccountNumber = account.Number
							break
						}
					}

					// get address
					handler.generatedAddress, handler.err = wallet.ReceiveAddress(handler.selectedAccountNumber)
					if handler.err != nil {
						contentWindow.SetErrorMessage(handler.err.Error())
					} else {
						window.Master().Changed()
					}
				}

				if handler.generatedAddress != "" {
					handler.RenderAddress(contentWindow)
				}
			}
			contentWindow.End()
		}
		pageWindow.End()
	}
}

func (handler *ReceiveHandler) RenderAddress(window *helpers.Window) {
	window.Row(50).Dynamic(1)
	window.LabelWrap("Address: " + handler.generatedAddress)

	// generate qrcode
	png, err := qrcode.New(handler.generatedAddress, qrcode.Medium)
	if err != nil {
		window.Row(300).Dynamic(1)
		window.LabelWrap(err.Error())
	} else {
		window.Row(200).Dynamic(1)
		img := png.Image(300)
		imgRGBA := image.NewRGBA(img.Bounds())
		draw.Draw(imgRGBA, img.Bounds(), img, image.Point{}, draw.Src)
		window.Image(imgRGBA)
	}
}
