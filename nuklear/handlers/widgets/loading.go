package widgets

import (
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

func ShowIsFetching(window *helpers.Window) {
	window.Row(30).Dynamic(1)
	window.Label("Fetching data...", "LC")
}
