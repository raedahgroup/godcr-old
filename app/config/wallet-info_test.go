package config

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet/util"
)

type walletDbDir struct {
	Source string
	Path   string
}

func detectWallets(ctx context.Context) ([]*WalletInfo, error) {
	var allDetectedWallets []*WalletInfo
	decredWalletDbDirectories := []walletDbDir{
		{Source: "dcrwallet", Path: dcrutil.AppDataDir("dcrwallet", false)},
		{Source: "decredition", Path: decreditionAppDirectory()},
		{Source: "godcr", Path: DefaultAppDataDir},
	}

	for _, walletDir := range decredWalletDbDirectories {
		detectedWallets, err := findWalletsInDirectory(walletDir.Path, walletDir.Source)
		if err != nil {
			return nil, fmt.Errorf("error searching for wallets: %s", err.Error())
		}
		allDetectedWallets = append(allDetectedWallets, detectedWallets...)
	}

	// mark default wallet
	if len(allDetectedWallets) == 1 {
		allDetectedWallets[0].Default = true
	}

	// update config file with detected wallets info
	err := UpdateConfigFile(func(config *ConfFileOptions) {
		config.Wallets = allDetectedWallets
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save %d detected wallets: %s", len(allDetectedWallets), err.Error())
	}

	return allDetectedWallets, nil
}

func findWalletsInDirectory(walletDir, walletSource string) (wallets []*WalletInfo, err error) {
	// netType checks if the name of the directory where a wallet.db file was found is the name of a known/supported network type
	// dcrwallet, decredition and dcrlibwallet place wallet db files in "mainnet" or "testnet3" directories
	// returns nil if the directory used does not correspond to a known/supported network type
	detectNetParams := func(path string) *netparams.Params {
		walletDbDir := filepath.Dir(path)
		netType := filepath.Base(walletDbDir)
		return util.NetParams(netType)
	}

	err = filepath.Walk(walletDir, func(path string, file os.FileInfo, err error) error {
		if err != nil || file.IsDir() || file.Name() != "wallet.db" {
			return nil
		}

		netParams := detectNetParams(path)
		if netParams == nil {
			return nil
		}

		wallets = append(wallets, &WalletInfo{
			DbDir:   filepath.Dir(path),
			Source:  walletSource,
			Network: netParams.Name,
		})
		return nil
	})
	return
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

func TestWalletInfo_MarshalFlag(t *testing.T) {
	var err error
	var wallets []*WalletInfo

	cfg, _, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	wallets = cfg.Wallets
	if wallets == nil {
		wallets, err = detectWallets(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	type fields struct {
		DbDir   string
		Network string
		Source  string
		Default bool
	}

	type test struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}

	tests := make([]test, len(wallets))
	wantFormat := `{"DbDir":"%s","Network":"%s","Source":"%s","Default":%t}`
	for i := range wallets {

		tests[i] = test{
			name: "wallet info marshal " + cfg.AppDataDir,
			fields: fields{
				DbDir:   wallets[i].DbDir,
				Network: wallets[i].Network,
				Source:  wallets[i].Source,
				Default: wallets[i].Default,
			},
			want:    fmt.Sprintf(wantFormat, wallets[i].DbDir, wallets[i].Network, wallets[i].Source, wallets[i].Default),
			wantErr: false,
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := &WalletInfo{
				DbDir:   test.fields.DbDir,
				Network: test.fields.Network,
				Source:  test.fields.Source,
				Default: test.fields.Default,
			}
			got, err := wallet.MarshalFlag()
			fmt.Println(got)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletInfo.MarshalFlag() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletInfo.MarshalFlag() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletInfo_UnmarshalFlag(t *testing.T) {
	var err error
	var wallets []*WalletInfo

	cfg, _, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	wallets = cfg.Wallets
	if wallets == nil {
		wallets, err = detectWallets(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	type fields struct {
		DbDir   string
		Network string
		Source  string
		Default bool
	}

	type test struct {
		name    string
		fields  fields
		value   string
		wantErr bool
	}

	valueFormat := `{"DbDir":"%s","Network":"%s","Source":"%s","Default":%t}`

	tests := make([]test, len(wallets))
	for i := range wallets {

		tests[i] = test{
			name: "wallet info unmarshal" + cfg.AppDataDir,
			fields: fields{
				DbDir:   wallets[i].DbDir,
				Network: wallets[i].Network,
				Source:  wallets[i].Source,
				Default: wallets[i].Default,
			},
			value:   fmt.Sprintf(valueFormat, wallets[i].DbDir, wallets[i].Network, wallets[i].Source, wallets[i].Default),
			wantErr: false,
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := &WalletInfo{
				DbDir:   test.fields.DbDir,
				Network: test.fields.Network,
				Source:  test.fields.Source,
				Default: test.fields.Default,
			}
			if err := wallet.UnmarshalFlag(test.value); (err != nil) != test.wantErr {
				t.Errorf("WalletInfo.UnmarshalFlag() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestWalletInfo_Summary(t *testing.T) {
	var err error
	var wallets []*WalletInfo

	cfg, _, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	wallets = cfg.Wallets
	if wallets == nil {
		wallets, err = detectWallets(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	type fields struct {
		DbDir   string
		Network string
		Source  string
		Default bool
	}
	type test struct {
		name   string
		fields fields
		want   string
	}

	tests := make([]test, len(wallets))
	for i := range wallets {
		tests[i] = test{
			name: "wallet summary " + cfg.AppDataDir,
			fields: fields{
				DbDir:   wallets[i].DbDir,
				Network: wallets[i].Network,
				Source:  wallets[i].Source,
				Default: wallets[i].Default,
			},
			want: fmt.Sprintf("%s wallet from %s", wallets[i].Network, wallets[i].Source),
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := &WalletInfo{
				DbDir:   test.fields.DbDir,
				Network: test.fields.Network,
				Source:  test.fields.Source,
				Default: test.fields.Default,
			}
			if got := wallet.Summary(); got != test.want {
				t.Errorf("WalletInfo.Summary() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDefaultWallet(t *testing.T) {
	var err error
	var wallets []*WalletInfo

	cfg, _, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	wallets = cfg.Wallets
	if wallets == nil {
		wallets, err = detectWallets(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	var defaultWallet *WalletInfo
	for i := range wallets {
		if wallets[i].Default {
			defaultWallet = wallets[i]
			break
		}
	}
	tests := []struct {
		name    string
		wallets []*WalletInfo
		want    *WalletInfo
	}{
		{
			name:    "default wallet",
			wallets: wallets,
			want:    defaultWallet,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := DefaultWallet(test.wallets); !reflect.DeepEqual(got, test.want) {
				t.Errorf("DefaultWallet() = %v, want %v", got, test.want)
			}
		})
	}
}
