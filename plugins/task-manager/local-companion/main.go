package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

const productionTaskManagerOrigin = "https://task-manager.xxsrez-work.chatgpt.site"
const uatTaskManagerOrigin = "https://task-manager-uat.xxsrez-work.chatgpt.site"

func main() {
	origin := flag.String("origin", productionTaskManagerOrigin, "Task Manager site origin")
	privateUATBridge := flag.Bool(
		"private-uat-bridge", false,
		"proxy the private UAT MCP through the local companion using a Keychain-held Sites bypass credential",
	)
	flag.Parse()
	if runtime.GOOS != "darwin" {
		_, _ = fmt.Fprintln(os.Stderr, "Task Manager local companion currently supports macOS only.")
		os.Exit(1)
	}
	httpClient := &http.Client{
		Timeout:   15 * time.Minute,
		Transport: http.DefaultTransport,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return errors.New("redirects are not accepted")
		},
	}
	if *privateUATBridge {
		transport, err := newSitesBypassTransport(
			*origin,
			httpClient.Transport,
			&keychainSitesBypassStore{account: uatTaskManagerOrigin},
		)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Task Manager local companion could not start: private UAT bridge configuration is invalid.")
			os.Exit(1)
		}
		httpClient.Transport = transport
	}
	manager, err := newOAuthManager(
		*origin, httpClient, &keychainCredentialStore{account: *origin}, systemBrowser{},
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Task Manager local companion could not start: invalid origin.")
		os.Exit(1)
	}
	uploader := newUploadClient(httpClient, *origin, manager)
	var remote remoteMCP
	if *privateUATBridge {
		remote = newRemoteMCPClient(httpClient, *origin, manager)
	}
	if err := serveMCP(context.Background(), os.Stdin, os.Stdout, uploader, remote); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Task Manager local companion stopped because stdio transport failed.")
		os.Exit(1)
	}
}
