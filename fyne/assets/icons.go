package assets

import (
	"github.com/gobuffalo/packr/v2"

	"fyne.io/fyne"
)

const (
	AccountsIcon = "account.png"
	HistoryIcon  = "history.png"
	OverviewIcon = "overview.png"
	ReceiveIcon  = "receive.png"
	SendIcon     = "send.png"
	StakingIcon  = "staking.png"
	DecredLogo   = "decred.png"

	Reveal    = "reveal.png"
	Conceal   = "conceal.png"
	Checkmark = "checkmark.png"
	Restore   = "restore.png"
	Add       = "add.png"
	Loader    = "loader.png"
	Back      = "back.png"
)

var imageBox = packr.New("icons", "../assets")

// Get returns a map from the names of the icons passed as arguments to
// the icon resources that correspond to them. If an error is encountered
// while loading any of the icons, the error is returned immediately.
func Get(names ...string) (map[string]*fyne.StaticResource, error) {

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
