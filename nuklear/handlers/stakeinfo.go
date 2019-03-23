package handlers

import (
	"context"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type StakeInfoHandler struct {
	err         error
	isRendering bool
	stakeInfo   *walletcore.StakeInfo
}

func (handler *StakeInfoHandler) BeforeRender() {
	handler.err = nil
	handler.stakeInfo = nil
	handler.isRendering = false
}

func (handler *StakeInfoHandler) Render(window *nucular.Window, wallet app.WalletMiddleware) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.stakeInfo, handler.err = wallet.StakeInfo(context.Background())
	}

	// draw page
	if pageWindow := helpers.NewWindow("StakeInfo Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("Stake Info")

		// content window
		if contentWindow := pageWindow.ContentWindow("StakeInfo Content"); contentWindow != nil {
			if handler.err != nil {
				contentWindow.SetErrorMessage(handler.err.Error())
			} else {
				contentWindow.Row(20).Static(43, 43, 43, 43, 43, 43, 80, 43, 43, 43, 60)
				contentWindow.Label("Expired", "LC")
				contentWindow.Label("Immature", "LC")
				contentWindow.Label("Live", "LC")
				contentWindow.Label("Revoked", "LC")
				contentWindow.Label("Unmined", "LC")
				contentWindow.Label("Unspent", "LC")
				contentWindow.Label("AllmempoolTix", "LC")
				contentWindow.Label("PoolSize", "LC")
				contentWindow.Label("Missed", "LC")
				contentWindow.Label("Voted", "LC")
				contentWindow.Label("Total Subsidy", "LC")

				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Expired)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Immature)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Live)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Revoked)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.OwnMempoolTix)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Unspent)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.AllMempoolTix)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.PoolSize)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Missed)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Voted)), "LC")
				contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.TotalSubsidy)), "LC")

			}
			contentWindow.End()
		}
		pageWindow.End()
	}
}
