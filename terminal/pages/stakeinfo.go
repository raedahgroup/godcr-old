package pages

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/rivo/tview"
)

func StakeinfoPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	stakeInfo, err := wallet.StakeInfo(context.Background())
	errmsg := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	if err != nil {
		return errmsg.SetText(fmt.Sprintf(err.Error()))
	}
	if stakeInfo == nil {
		return errmsg.SetText(fmt.Sprintf("no tickets in wallet"))
	}

	body := tview.NewTable().SetBorders(true)
	body.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			body.ScrollToBeginning()
			clearFocus()
		}
	})
	body.SetCell(0, 0, tview.NewTableCell("Expired").SetAlign(tview.AlignCenter))
	body.SetCell(0, 1, tview.NewTableCell("Immature").SetAlign(tview.AlignCenter))
	body.SetCell(0, 2, tview.NewTableCell("Live").SetAlign(tview.AlignCenter))
	body.SetCell(0, 3, tview.NewTableCell("Revoked").SetAlign(tview.AlignCenter))
	body.SetCell(0, 4, tview.NewTableCell("Unmined").SetAlign(tview.AlignCenter))
	body.SetCell(0, 5, tview.NewTableCell("Unspent").SetAlign(tview.AlignCenter))
	body.SetCell(0, 6, tview.NewTableCell("AllmempoolTix").SetAlign(tview.AlignCenter))
	body.SetCell(0, 7, tview.NewTableCell("PoolSize").SetAlign(tview.AlignCenter))
	body.SetCell(0, 8, tview.NewTableCell("Missed").SetAlign(tview.AlignCenter))
	body.SetCell(0, 9, tview.NewTableCell("Voted").SetAlign(tview.AlignCenter))
	body.SetCell(0, 10, tview.NewTableCell("Total Subsidy").SetAlign(tview.AlignCenter))

	body.SetCell(1, 0, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Expired))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 1, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Immature))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 2, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Live))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 3, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Revoked))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 4, tview.NewTableCell(strconv.Itoa(int(stakeInfo.OwnMempoolTix))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 5, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Unspent))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 6, tview.NewTableCell(strconv.Itoa(int(stakeInfo.AllMempoolTix))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 7, tview.NewTableCell(strconv.Itoa(int(stakeInfo.PoolSize))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 8, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Missed))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 9, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Voted))).SetAlign(tview.AlignCenter))
	body.SetCell(1, 10, tview.NewTableCell(strconv.Itoa(int(stakeInfo.TotalSubsidy))).SetAlign(tview.AlignCenter))

	setFocus(body)

	return body
}
