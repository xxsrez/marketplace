//go:build darwin

package main

import (
	"bytes"
	"os/exec"
	"strings"
)

const sitesBypassKeychainService = "com.xxsrez.task-manager.uat.sites-bypass"

type keychainSitesBypassStore struct {
	account string
}

func (store *keychainSitesBypassStore) LoadSecret() (string, error) {
	command := exec.Command(
		"/usr/bin/security", "find-generic-password", "-a", store.account,
		"-s", sitesBypassKeychainService, "-w",
	)
	var stdout bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = nil
	if err := command.Run(); err != nil {
		return "", newLocalError(
			"sites_bypass_unavailable",
			"private UAT access credential is unavailable in macOS Keychain",
		)
	}
	secret := strings.TrimSpace(stdout.String())
	if secret == "" || strings.ContainsAny(secret, "\r\n") {
		return "", newLocalError("sites_bypass_invalid", "private UAT access credential is invalid")
	}
	return secret, nil
}
