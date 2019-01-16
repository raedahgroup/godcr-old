package pages

import (
	"fmt"
	"github.com/raedahgroup/godcr/app"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type statusPage struct{}

func (s *statusPage) Setup() *widgets.QWidget {
	pageContent := widgets.NewQWidget(nil, 0)

	// create layout to arrange child views vertically and center them
	pageLayout := widgets.NewQVBoxLayout()
	pageLayout.SetAlign(core.Qt__AlignCenter)
	pageContent.SetLayout(pageLayout)

	// add views to page layout
	statusLabel := widgets.NewQLabel2(fmt.Sprintf("%s status: running", app.Name), nil, 0)
	pageContent.Layout().AddWidget(statusLabel)

	return pageContent
}
