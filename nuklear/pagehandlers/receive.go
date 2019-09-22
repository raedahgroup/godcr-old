package pagehandlers

import (
	"image"
	"image/draw"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/atotto/clipboard"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/nuklog"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
	"github.com/skip2/go-qrcode"
)

const (
	qrCodeImageSize                      = 300
	qrCodeAddressHolderHorizontalPadding = 40
)

type ReceiveHandler struct {
	wallet                walletcore.Wallet
	accountSelectorWidget *widgets.AccountSelector
	generateAddressError  error
	generatedAddress      string
	accounts              []*walletcore.Account
	accountNumber         uint32
}

func (handler *ReceiveHandler) BeforeRender(wallet walletcore.Wallet, settings *config.Settings, refreshWindowDisplay func()) bool {
	handler.wallet = wallet
	handler.generateAddressError = nil
	handler.generatedAddress = ""

	handler.accountSelectorWidget = widgets.AccountSelectorWidget("Account:", false, false, wallet, nil)
	handler.accounts, handler.generateAddressError = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	
	handler.accountNumber = handler.accounts[0].Number
	handler.generatedAddress, handler.generateAddressError = handler.wallet.ReceiveAddress(handler.accountNumber)
	
	if len(handler.accounts) > 1 {
		handler.accountSelectorWidget = widgets.AccountSelectorWidget("Account:", false, false, wallet, func(int){
			handler.accountNumber = handler.accountSelectorWidget.GetSelectedAccountNumber()
			handler.generatedAddress, handler.generateAddressError = handler.wallet.ReceiveAddress(handler.accountNumber)
		})
	}
	return true
}

func (handler *ReceiveHandler) Render(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("Receiving Decred", window, func(contentWindow *widgets.Window) {
		contentWindow.AddWrappedLabelWithColor(walletcore.ReceivingDecredHint, widgets.LeftCenterAlign, styles.GrayColor)

		contentWindow.AddHorizontalSpace(1)

		contentWindow.Row(widgets.EditorHeight).Static(contentWindow.LabelWidth("You can also manually generate a"), 120)
		contentWindow.AddLabelsToCurrentRow(widgets.NewColoredLabelTableCell("You can also manually generate a", widgets.LeftCenterAlign, styles.GrayColor))
		contentWindow.AddLinkLabelsToCurrentRow(widgets.NewLinkLabelCellCell("new address.", func() {
			handler.generatedAddress, handler.generateAddressError = handler.wallet.GenerateNewAddress(handler.accountNumber)
			window.Master().Changed()
		}))

		contentWindow.AddHorizontalSpace(15)

		// draw account selection widget before rendering previously generated address
		handler.accountSelectorWidget.Render(contentWindow)
		window.Master().Changed()

		// display error if there was an error the last time address generation was attempted
		if handler.generateAddressError != nil {
			contentWindow.DisplayErrorMessage("Address could not be generated", handler.generateAddressError)
		} else if handler.generatedAddress != "" {
			handler.RenderAddress(contentWindow)
		}
	})
}

func (handler *ReceiveHandler) RenderAddress(window *widgets.Window) {
	generatedAddressWidth := window.LabelWidth(handler.generatedAddress)
	qrCodeAddressHolderWidth, qrCodeAddressHolderHeight := qrCodeImageSize, qrCodeImageSize
	if generatedAddressWidth >= qrCodeImageSize {
		qrCodeAddressHolderWidth = generatedAddressWidth
	}
	qrCodeAddressHolderWidth += qrCodeAddressHolderHorizontalPadding
	qrCodeAddressHolderHeight += window.SingleLineLabelHeight()

	// generate qrcode
	qrCode, err := qrcode.New(handler.generatedAddress, qrcode.Medium)
	if err != nil {
		// todo logs need to accept message to accompany errors
		nuklog.LogError(err)
		window.DisplayErrorMessage("Error generating qr code", err)
		window.AddLabel(handler.generatedAddress, widgets.LeftCenterAlign)
	} else {
		sourceImage := qrCode.Image(qrCodeImageSize)
		qrCodeImage := image.NewRGBA(sourceImage.Bounds())
		draw.Draw(qrCodeImage, sourceImage.Bounds(), sourceImage, image.Point{}, 0)

		// holder for code and address
		window.Row(qrCodeAddressHolderHeight).SpaceBegin(2)

		// calculate left padding space to use before displaying image to place in horizontal center
		qrCodeImageLeftPadding := (qrCodeAddressHolderWidth - qrCodeImageSize) / 2
		window.LayoutSpacePushScaled(rect.Rect{
			X: qrCodeImageLeftPadding,
			Y: 0,
			W: qrCodeImageSize,
			H: qrCodeImageSize,
		})
		window.Image(qrCodeImage)

		// calculate left padding space to use before displaying address label to place in horizontal center
		addressLabelLeftPadding := (qrCodeAddressHolderWidth - generatedAddressWidth) / 2
		window.LayoutSpacePushScaled(rect.Rect{
			X: addressLabelLeftPadding,
			Y: qrCodeImageSize,
			W: generatedAddressWidth,
			H: window.SingleLineLabelHeight(),
		})
		var addressClicked bool
		window.SelectableLabel(handler.generatedAddress, widgets.CenterAlign, &addressClicked)

		if addressClicked {
			clipboard.WriteAll(handler.generatedAddress)
		}
	}
}
