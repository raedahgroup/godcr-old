package stakingpagehandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
)

func (stakingPage *StakingPageObjects) initBaseObjects() error {
	icons, err := assets.GetIcons(assets.CollapseIcon)
	if err != nil {
		return err
	}
	stakingPage.icons = icons

	stakingPageLabel := widget.NewLabelWithStyle(values.StakingPageTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	stakingPage.StakingPageContents.Append(widget.NewHBox(stakingPageLabel))
	return nil
}
