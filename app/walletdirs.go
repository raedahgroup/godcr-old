package app

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/decred/dcrd/dcrutil"
)

type WalletDbDir struct {
	Source string
	Path   string
}

// WalletDbFileName is the name used by dcrwallet, decredition and dcrlibwallet when creating wallets
const WalletDbFileName = "wallet.db"

// DecredWalletDbDirectories maintains a slice of directories where decred wallet databases may be found
func DecredWalletDbDirectories() (directories []WalletDbDir) {
	// scan for all potential dcrwallet directories and return
	dcrWalletWildCardDir := dcrutil.AppDataDir("dcrwallet*", false)
	dcrWalletDirs, _ := filepath.Glob(dcrWalletWildCardDir)
	for _, dcrWalletDir := range dcrWalletDirs {
		directories = append(directories, WalletDbDir{
			Source: "dcrwallet",
			Path:   dcrWalletDir,
		})
	}

	// scan for all potential decredition directories and return
	decreditionWildCardDir := decreditionAppDirectory("*")
	decreditionWalletDirs, _ := filepath.Glob(decreditionWildCardDir)
	for _, decreditionWalletDir := range decreditionWalletDirs {
		directories = append(directories, WalletDbDir{
			Source: "decredition",
			Path:   decreditionWalletDir,
		})
	}

	// scan for all potential godcr directories and return
	godcrWildCardDir := dcrutil.AppDataDir("godcr*", false)
	godcrWalletDirs, _ := filepath.Glob(godcrWildCardDir)
	for _, godcrWalletDir := range godcrWalletDirs {
		directories = append(directories, WalletDbDir{
			Source: "godcr",
			Path:   godcrWalletDir,
		})
	}

	return
}

// decreditionAppDirectory returns the appdata dir used by decredition on different operating systems
// following the pattern in the decredition source code
// see https://github.com/decred/decrediton/blob/master/app/main_dev/paths.js#L10-L18
func decreditionAppDirectory(wildcard string) string {
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
		return filepath.Join(homeDir, "AppData", "Local", "Decrediton"+wildcard)
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "decrediton"+wildcard)
	default:
		return filepath.Join(homeDir, ".config", "decrediton"+wildcard)
	}
}
