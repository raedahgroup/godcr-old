package createwallet


import (
	"gioui.org/layout"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

type CreateWalletHandler struct {
}

func NewCreateWalletHandler() *CreateWalletHandler {
	return &CreateWalletHandler{}
}

func (h *CreateWalletHandler) BeforeRender(walletMiddleware app.WalletMiddleware, config *config.Settings) {

}

func (h *CreateWalletHandler) Render(ctx *layout.Context, refreshWindowFunc func()) {

}