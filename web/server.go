package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
	"github.com/raedahgroup/godcr/web/routes"
	"github.com/raedahgroup/godcr/web/weblog"
)

func StartServer(ctx context.Context, walletMiddleware app.WalletMiddleware, httpHost, httpPort string, settings *config.Settings) error {
	router := chi.NewRouter()

	// setup static file serving
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "../../web/static/dist")
	makeStaticFileServer(router, "/static", http.Dir(filesDir))

	// setup routes for templated pages, returns sync blockchain function if wallet is successfully opened
	// returns error if wallet exists but could not be opened
	syncBlockchain, err := routes.OpenWalletAndSetupRoutes(ctx, walletMiddleware, router, settings)
	if err != nil {
		return err
	}

	fmt.Println("Starting web server")

	serverAddress := net.JoinHostPort(httpHost, httpPort)
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
		go askToLaunchBrowser(address) // run in goroutine so this function returns immediately without waiting for user response
		return nil
	}
}

func askToLaunchBrowser(address string) {
	launchBrowserConfirmed, err := terminalprompt.RequestYesNoConfirmation("Do you want to launch the web browser?", "")
	if err != nil {
		weblog.Log.Error("Failed to read input", err.Error())
		fmt.Fprintf(os.Stderr, "Error reading your response: %s\n", err.Error())
		return
	}

	if !launchBrowserConfirmed {
		return
	}

	fmt.Print("Launching browser... ") // use print so next text can be added to same line

	if launchError := launchBrowser("http://" + address); launchError != nil {
		weblog.Log.Error("Failed to launch browser", launchError.Error())
		fmt.Println("Browser failed to launch.")
	} else {
		fmt.Println("Done.")
	}
}

func launchBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}
