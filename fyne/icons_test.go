package fyne

import (
	"testing"
)

func TestGetIcons(t *testing.T) {
	t.Run("given icons that exist in the assets directory, maps each one to a fyne static resource", func(t *testing.T) {
		iconNames := []string{
			aboutIcon, accountsIcon, decredDarkIcon, decredLightIcon, exitIcon, helpIcon, historyIcon, moreIcon,
			overviewIcon, receiveIcon, securityIcon, sendIcon, settingsIcon, stakeIcon}
		icons, err := getIcons(iconNames...)
		if err != nil {
			t.Errorf("getIcons returned an error for icons in the assets directory: %v", err)
		}
		about, ok := icons[aboutIcon]
		if !ok {
			t.Errorf("could not find %s in the icons map. Check that it is in the assets directory", aboutIcon)
		}
		if about.Name() != aboutIcon {
			t.Errorf("")
			t.Errorf("unexpected static resource found: expected %s, got %s", aboutIcon, about.Name())
		}
	})
	t.Run("returns an error for icons that do not exist in the assets directory", func(t *testing.T) {
		iconNames := []string{"not_existing.png", "deleted.png"}
		_, err := getIcons(iconNames...)
		if err == nil {
			t.Errorf("expected an error for icon names %v, got nil", iconNames)
		}
	})
}
