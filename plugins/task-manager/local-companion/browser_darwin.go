//go:build darwin

package main

import "os/exec"

type systemBrowser struct{}

func (systemBrowser) Open(rawURL string) error {
	return exec.Command("/usr/bin/open", rawURL).Run()
}
