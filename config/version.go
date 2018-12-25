package config

import (
	"fmt"
)

// version provides version information for the program.
type version struct {
	major, minor, patch int
	label               string
	nick                string
}

var ver = version{
	major: 0,
	minor: 0,
	patch: 1,
	label: "",
}

// CommitHash may be set on the build command line:
// go build -ldflags "-X github.com/decred/dcrdata/version.CommitHash=`git describe --abbrev=8 --long | awk -F "-" '{print $(NF-1)"-"$NF}'`"
var CommitHash string

func (v *version) String() string {
	var hashStr string
	if CommitHash != "" {
		hashStr = "+" + CommitHash
	}
	if v.label != "" {
		return fmt.Sprintf("%d.%d.%d-%s%s",
			v.major, v.minor, v.patch, v.label, hashStr)
	}
	return fmt.Sprintf("%d.%d.%d%s",
		v.major, v.minor, v.patch, hashStr)
}

// AppVersion provides the version string for the program.
func AppVersion() string {
	return fmt.Sprintf("%s version: %s", AppName(), ver.String())
}
