package handlers

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type SendHandler struct {
	err          error
	hasRendered  bool
	transactions []*walletcore.Transaction
	wallet walletcore.Wallet
}

func (handler *SendHandler) SetWalletMiddleware(walletMiddleare walletcore.Wallet) {
	handler.wallet = walletMiddleare
}

func (handler *SendHandler) BeforeRender() {
	handler.hasRendered = true
}

func (handler *SendHandler) Render(w *nucular.Window) {

}

