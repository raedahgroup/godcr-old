package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Name returns the name of the binary file that started this program
func Name() string {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	return appName
}

// version provides version information for the program
type version struct {
	major, minor, patch int
	label               string
}

var currentVersion = version{
	major: 0,
	minor: 0,
	patch: 1,
	label: "",
}

// todo this comment needs correction
// CommitHash may be set on the build command line:
// go build -ldflags "-X github.com/decred/dcrdata/version.CommitHash=`git describe --abbrev=8 --long | awk -F "-" '{print $(NF-1)"-"$NF}'`"
var CommitHash string

// Version returns the version of this app in a easy-to-read format
func Version() string {
	var hashStr string
	if CommitHash != "" {
		hashStr = "+" + CommitHash
	}

	if currentVersion.label == "" {
		return fmt.Sprintf("%d.%d.%d%s", currentVersion.major, currentVersion.minor, currentVersion.patch, hashStr)
	}

	return fmt.Sprintf("%d.%d.%d-%s%s", currentVersion.major, currentVersion.minor, currentVersion.patch, currentVersion.label, hashStr)
}
