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

func main() {
	origin := flag.String("origin", productionTaskManagerOrigin, "Task Manager site origin")
	transportOrigin := flag.String(
		"transport-origin", "", "Optional exact 127.0.0.1 ingress used for machine requests",
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
	credentials := &volatileCredentialStore{}
	manager, err := newOAuthManager(
		*origin, *transportOrigin, httpClient, credentials, systemBrowser{},
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Task Manager local companion could not start: invalid origin.")
		os.Exit(1)
	}
	requestOrigin := *transportOrigin
	if requestOrigin == "" {
		requestOrigin = *origin
	}
	service := newUploadClient(httpClient, requestOrigin, manager)
	if err := serveMCP(context.Background(), os.Stdin, os.Stdout, service); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Task Manager local companion stopped because stdio transport failed.")
		os.Exit(1)
	}
}
