package handlers

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type SendHandler struct {
	err          error
	hasRendered  bool
	transactions []*walletcore.Transaction
}

func (handler *SendHandler) BeforeRender() {
	handler.hasRendered = true
}

func (handler *SendHandler) Render(w *nucular.Window, walletMiddleware app.WalletMiddleware) {

}
