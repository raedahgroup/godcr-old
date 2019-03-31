package handlers

import (
	"image"
	"image/draw"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
	qrcode "github.com/skip2/go-qrcode"
)

type ReceiveHandler struct {
	err         error
	isRendering bool

	numAccounts     int
	accountSelector *widgets.AccountSelection

	addressInput nucular.TextEditor

	// generatedAddress
	generatedAddress string
}

func (handler *ReceiveHandler) BeforeRender() {
	handler.err = nil
	handler.isRendering = false

	handler.addressInput.Flags = nucular.EditClipboard | nucular.EditNoCursor
}

func (handler *ReceiveHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		accounts, err := walletMiddleware.AccountsOverview(walletcore.DefaultRequiredConfirmations)
		if err != nil {
			handler.err = err
		} else {
			handler.numAccounts = len(accounts)
			handler.accountSelector = widgets.NewAccountSelectionWidget(accounts)
		}
	}

	// draw page
	if pageWindow := helpers.NewWindow("Receive Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("Generate Receive Address")

		// content window
		if contentWindow := pageWindow.ContentWindow("Receive Content"); contentWindow != nil {
			if handler.err != nil {
				contentWindow.SetErrorMessage(handler.err.Error())
			} else {
				contentWindow.Row(30).Ratio(0.75, 0.25)

				buttonLabel := "Generate"
				if handler.numAccounts == 1 {
					buttonLabel = "Regenerate"
				}

				if handler.numAccounts == 1 && handler.generatedAddress == "" {
					handler.generatedAddress, handler.err = walletMiddleware.ReceiveAddress(handler.accountSelector.GetSelectedAccountNumber())
					handler.RenderAddress(contentWindow)
				} else {
					handler.accountSelector.Render(contentWindow.Window)

					// draw submit button
					if contentWindow.Button(label.T(buttonLabel), false) {
						// get address
						handler.generatedAddress, handler.err = walletMiddleware.ReceiveAddress(handler.accountSelector.GetSelectedAccountNumber())
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

			}
			contentWindow.End()
		}
		pageWindow.End()
	}
}

func (handler *ReceiveHandler) RenderAddress(window *helpers.Window) {
	window.Row(30).Static(40, 400)
	window.Label("Address: ", "LC")
	helpers.StyleClipboardInput(window)
	handler.addressInput.Buffer = []rune(handler.generatedAddress)
	handler.addressInput.Edit(window.Window)
	helpers.ResetInputStyle(window)

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
