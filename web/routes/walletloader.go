package routes

import (
	"fmt"
	"net/http"
)

func (routes *Routes) makeWalletLoaderMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return routes.walletLoaderFn(next)
	}
}

// walletLoaderFn checks if wallet is not open and attempts to open it
// if an error occurs while attempting to open wallet, an error page is displayed and the actual route handler is not called
func (routes *Routes) walletLoaderFn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !routes.walletMiddleware.IsWalletOpen() {
			err := routes.loadWallet()
			if err != nil {
				routes.renderError(err.Error(), res)
				return
			}
		}

		next.ServeHTTP(res, req)
	})
}

func (routes *Routes) loadWallet() error {
	walletExists, err := routes.walletMiddleware.WalletExists()
	if err != nil {
		return fmt.Errorf("Error checking for wallet: %s", err.Error())
	}

	if !walletExists {
		return fmt.Errorf("Wallet not created. Please create a wallet to continue. Use `dcrcli create` on terminal")
	}

	err = routes.walletMiddleware.OpenWallet()
	if err != nil {
		return fmt.Errorf("Failed to open wallet: %s", err.Error())
	}

	return nil
}
