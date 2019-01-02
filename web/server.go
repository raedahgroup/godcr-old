package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/web/routes"
)

func StartHttpServer(ctx context.Context, walletMiddleware app.WalletMiddleware, address string) error {
	router := chi.NewRouter()

	// first try to load wallet if it exists
	err := openWalletIfExist(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	// setup static file serving
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web/public")
	makeStaticFileServer(router, "/static", http.Dir(filesDir))

	// setup routes for templated pages, returns wallet loader function
	syncBlockchain := routes.Setup(walletMiddleware, router)

	fmt.Println("Starting web server")
	err = startServer(ctx, address, router)
	if err != nil {
		return err
	}

	// check if context has been canceled before starting blockchain sync
	err = ctx.Err()
	if err != nil {
		fmt.Println("Web server stopped")
		return err
	}
	syncBlockchain()

	// keep alive till ctx is canceled
	<-ctx.Done()
	fmt.Println("Web server stopped")
	return nil
}

// this method may stall until previous godcr instances are closed (especially in cases of multiple dcrlibwallet instances)
// hence the need for ctx, so user can cancel the operation if it's taking too long
func openWalletIfExist(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	// notify user of the current operation so if takes too long, they have an idea what the cause is
	fmt.Println("Opening wallet...")

	var err error
	var errMsg string
	loadWalletDone := make(chan bool)

	go func() {
		defer func() {
			loadWalletDone <- true
		}()

		var walletExists bool
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
		return err

	case <-ctx.Done():
		return ctx.Err()
	}
}

func makeStaticFileServer(router chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		router.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	router.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// startServer attempts to listen for connections on `address` in a goroutine
// any error that occur during the attempt are broadcasted to `errChan`
// startServer waits 2 seconds to catch error sent to `errChan` and returns the error
// startServer returns nil, if no error was received during the 2-seconds window
// startServer returns error if ctx is canceled while waiting
func startServer(ctx context.Context, address string, router chi.Router) error {
	// check if context has been canceled before attempting to start server
	err := ctx.Err()
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(address, router)
	}()

	// briefly wait for an error and then return
	t := time.NewTimer(2 * time.Second)
	select {
	case err := <-errChan:
		fmt.Fprintf(os.Stderr, "Web server failed to start: %s\n", err.Error())
		return err
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "Web server not started")
		return ctx.Err()
	case <-t.C:
		fmt.Printf("Web server running on %s\n", address)
		return nil
	}
}
