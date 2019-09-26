package pages

import "github.com/rivo/tview"

func ExitPage() tview.Primitive {
	body := tview.NewModal().
		SetText("Do you want to quit Terminal application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				commonPageData.app.Stop()
			} else {
				commonPageData.clearAllPageContent()
			}
		})

	commonPageData.app.SetFocus(body)
	return body
}
