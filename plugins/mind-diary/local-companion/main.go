package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

const mindDiaryUATOrigin = "https://mind-diary.xxsrez-work.chatgpt.site"

func main() {
	if runtime.GOOS != "darwin" {
		_, _ = fmt.Fprintln(os.Stderr, "Mind Diary local companion currently supports macOS only.")
		os.Exit(1)
	}
	client := &http.Client{
		Timeout: 15 * time.Minute,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return errors.New("redirects are not accepted")
		},
	}
	workspaceRoots, err := configuredWorkspaceRoots()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Mind Diary local companion could not start: workspace authority is invalid.")
		os.Exit(1)
	}
	service, err := newLocalFileService(client, mindDiaryUATOrigin, workspaceRoots...)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Mind Diary local companion could not start.")
		os.Exit(1)
	}
	defer service.Close()
	if err := serveMCP(context.Background(), os.Stdin, os.Stdout, service); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Mind Diary local companion stopped because stdio transport failed.")
		os.Exit(1)
	}
}
