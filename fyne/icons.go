package fyne

import (
	"github.com/gobuffalo/packr/v2"

	"fyne.io/fyne"
)

const (
	aboutIcon       = "about.png"
	accountsIcon    = "account.png"
	decredDarkIcon  = "decredDark.png"
	decredLightIcon = "decredLight.png"
	exitIcon        = "exit.png"
	helpIcon        = "help.png"
	historyIcon     = "history.png"
	moreIcon        = "more.png"
	overviewIcon    = "overview.png"
	receiveIcon     = "receive.png"
	securityIcon    = "security.png"
	sendIcon        = "send.png"
	settingsIcon    = "settings.png"
	stakeIcon       = "stake.png"
)

var imageBox = packr.New("icons", "assets")

// getIcons returns a map from the names of the icons passed as arguments to
// the icon resources that correspond to them. If an error is encountered
// while loading any of the icons, the error is returned immediately.
func getIcons(names ...string) (map[string]*fyne.StaticResource, error) {
	icons := make(map[string]*fyne.StaticResource, len(names))
	for _, name := range names {
		iconBytes, err := imageBox.Find(name)
		if err != nil {
			return nil, err
		}
		icons[name] = &fyne.StaticResource{StaticName: name, StaticContent: iconBytes}
	}
	return icons, nil
}
