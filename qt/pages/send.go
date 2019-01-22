package pages

import (
	"context"
	//"fmt"
	//"strings"

	//"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type sendPage struct {
	pageStub
	accountSourceDisplayLabel *widgets.QLabel
	accountNameInput          *widgets.QLineEdit
	spendUnconfirmedLabel     *widgets.QLabel
	amountTitleLabel          *widgets.QLabel
	amountInput               *widgets.QLabel
	destinationAddressLabel   *widgets.QLabel
	destinantionAddressInput  *widgets.QFormLayout
	selectCustomInput         *widgets.QLabel
	sendLayout                *widgets.QGridLayout
}

func (s *sendPage) SetupWithWallet(ctx context.Context, wallet walletcore.Wallet) *widgets.QWidget {

	pageContent := widgets.NewQWidget(nil, 0)

	pageLayout := widgets.NewQVBoxLayout()
	pageLayout.SetAlign(core.Qt__AlignTop)
	pageContent.SetLayout(pageLayout)

	s.accountSourceDisplayLabel = widgets.NewQLabel2("Source Account", nil, 0)
	s.accountNameInput = widgets.NewQLineEdit(nil)
	s.accountNameInput.SetPlaceholderText("Enter Account Number")

	pageContent.Layout().AddWidget(s.accountSourceDisplayLabel)
	pageContent.Layout().AddWidget(s.accountNameInput)

	return pageContent

}

func (b *sendPage) sendToken(wallet walletcore.Wallet) {

}
