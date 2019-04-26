package pages

import (
	"context"

	"github.com/raedahgroup/godcr/app/walletcore"
)

type OverviewHandler struct {
	ctx    context.Context
	wallet walletcore.Wallet
}

func (handler *OverviewHandler) BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool {
	return false
}

func (handler *OverviewHandler) Render() {

}
