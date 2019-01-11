package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/decred/dcrd/dcrutil"
)

type walletDbDir struct {
	Source string
	Path string
}

// WalletDbFileName is the name used by dcrwallet, decredition and dcrlibwallet when creating wallets
const WalletDbFileName = "wallet.db"

// DecredWalletDbDirectories maintains a slice of directories where decred wallet databases may be found
func DecredWalletDbDirectories() []walletDbDir {
	return []walletDbDir{
		{ Source: "dcrwallet", Path: dcrutil.AppDataDir("dcrwallet", false) },
		{ Source: "decredition", Path: decreditionAppDirectory() },
		{ Source: "godcr", Path: defaultAppDataDir },
	}
}

// decreditionAppDirectory returns the appdata dir used by decredition on different operating systems
// following the pattern in the decredition source code
// see https://github.com/decred/decrediton/blob/master/app/main_dev/paths.js#L10-L18
func decreditionAppDirectory() string {
	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "AppData", "Local", "Decrediton")
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "decrediton")
	default:
		return filepath.Join(homeDir, ".config", "decrediton")
	}
}
