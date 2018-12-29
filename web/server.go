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
	"github.com/raedahgroup/dcrcli/app"
	"github.com/raedahgroup/dcrcli/web/routes"
)

func StartHttpServer(ctx context.Context, walletMiddleware app.WalletMiddleware, address string) {
	router := chi.NewRouter()

	// setup static file serving
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web/public")
	makeStaticFileServer(router, "/static", http.Dir(filesDir))

	// setup routes for templated pages, returns wallet loader function
	loadWalletAndSyncBlockchain := routes.Setup(walletMiddleware, router)

	fmt.Println("Starting web server")
	err := startServer(ctx, address, router)
	if err != nil {
		return
	}

	// check if context has been canceled before attempting to load wallet
	err = ctx.Err()
	if err != nil {
		fmt.Println("Web server stopped")
		return
	}
	go loadWalletAndSyncBlockchain()

	// keep alive till ctx is canceled
	<-ctx.Done()
	fmt.Println("Web server stopped")
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
	case <- ctx.Done():
		fmt.Fprintln(os.Stderr, "Web server not started")
		return ctx.Err()
	case <-t.C:
		fmt.Printf("Web server running on %s\n", address)
		return nil
	}
}
