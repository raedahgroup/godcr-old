package terminal

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raedahgroup/godcr/app"
)

// this method may stall until previous godcr instances are closed (especially in cases of multiple dcrlibwallet instances)
// hence the need for ctx, so user can cancel the operation if it's taking too long
func openWalletIfExist(ctx context.Context, walletMiddleware app.WalletMiddleware) (walletExists bool, err error) {
	var errMsg string
	loadWalletDone := make(chan bool)

	go func() {
		defer func() {
			loadWalletDone <- true
		}()

		walletExists, err = walletMiddleware.WalletExists()
		if err != nil {
			errMsg = fmt.Sprintf("Error checking %s wallet", walletMiddleware.NetType())
		}
		if err != nil || !walletExists {
			return
		}

		err = walletMiddleware.OpenWallet()
		if err != nil {
			errMsg = fmt.Sprintf("Failed to open %s wallet", walletMiddleware.NetType())
		}
	}()

	select {
	case <-loadWalletDone:
		if errMsg != "" {
			fmt.Fprintln(os.Stderr, errMsg)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return

	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func CreateWallet(ctx context.Context, password string, walletMiddleware app.WalletMiddleware) (err error) {
	// first check if wallet already exists
	//work in progress
	walletExists, err := walletMiddleware.WalletExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking %s wallet: %s\n", walletMiddleware.NetType(), err.Error())
		return
	}
	if walletExists {
		netType := strings.Title(walletMiddleware.NetType())
		fmt.Fprintf(os.Stderr, "%s wallet already exists\n", netType)
		return fmt.Errorf("wallet already exists")
	}

	// get seed and display to user
	seed, err := walletMiddleware.GenerateNewWalletSeed()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating seed for new wallet: %s\n", err)
		return
	}
	//todo show user seed that was generated
	err = walletMiddleware.CreateWallet(password, seed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating wallet: %s\n", err.Error())
		return
	}
	fmt.Printf("Decred %s wallet created successfully\n", walletMiddleware.NetType())

	// todo sync blockchain?
	return nil
}
