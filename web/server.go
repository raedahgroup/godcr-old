package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/browser"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/web/routes"
)

func StartServer(ctx context.Context, walletMiddleware app.WalletMiddleware, host, port string) error {
	router := chi.NewRouter()

	// setup static file serving
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web/static/dist")
	makeStaticFileServer(router, "/static", http.Dir(filesDir))

	// setup routes for templated pages, returns sync blockchain function if wallet is successfully opened
	// returns error if wallet exists but could not be opened
	syncBlockchain, err := routes.OpenWalletAndSetupRoutes(ctx, walletMiddleware, router)
	if err != nil {
		return err
	}

	fmt.Println("Starting web server")

	serverAddress := net.JoinHostPort(host, port)
	err = startServer(ctx, serverAddress, router)
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
		fmt.Printf("Launching browser....\n")
		browser.OpenURL("http://" + address)
		return nil
	}
}
