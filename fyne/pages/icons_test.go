package pages

import (
	"math/rand"
	"testing"
)

func TestGetIcons(t *testing.T) {
	t.Run("given icons that exist in the assets directory, maps each one to a fyne static resource", func(t *testing.T) {
		iconNames := []string{
			accountsIcon, historyIcon, overviewIcon, receiveIcon, sendIcon, stakeIcon}
		icons, err := getIcons(iconNames...)
		if err != nil {
			t.Errorf("getIcons returned an error for icons in the assets directory: %v", err)
			t.FailNow()
		}
		randomIconName := iconNames[rand.Intn(len(iconNames))]
		randomIcon, ok := icons[randomIconName]
		if !ok {
			t.Errorf("could not find %s in the icons map. Check that it is in the assets directory", randomIconName)
		}
		if randomIcon.Name() != randomIconName {
			t.Errorf("")
			t.Errorf("unexpected static resource found: expected %s, got %s", randomIconName, randomIcon.Name())
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
