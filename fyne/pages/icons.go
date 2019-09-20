package pages

import (
	"github.com/gobuffalo/packr/v2"

	"fyne.io/fyne"
)

const (
	accountsIcon    = "account.png"
	historyIcon     = "history.png"
	overviewIcon    = "overview.png"
	receiveIcon     = "receive.png"
	sendIcon        = "send.png"
	stakeIcon       = "stake.png"
	decredLogo      = "decred.png"
	reveal          = "reveal.png"
	conceal         = "conceal.png"
	checkmark       = "checkmark.png"
	createNewWallet = "createNewWallet.png"
	restoreWallet   = "restoreWallet.png"
	collapse        = "collapse.png"
	wordlist        = "wordlist.txt"
)

var imageBox = packr.New("icons", "../assets")

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
