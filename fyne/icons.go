package fyne

import (
	"fmt"
	"path/filepath"
	"sync"

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

var (
	imageBox       *packr.Box
	prepareBoxOnce = sync.Once{}
)

func makeIconBox() (*packr.Box, error) {
	var iconsLocation, err = filepath.Abs("assets")
	if err != nil {
		return nil, fmt.Errorf("could not get path to assets directory: %v", err)
	}
	return packr.New("icons", iconsLocation), nil
}

func getIcons(names ...string) (map[string]*fyne.StaticResource, error) {
	var err error
	prepareBoxOnce.Do(func() {
		imageBox, err = makeIconBox()
	})
	if err != nil {
		return nil, err
	}
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
