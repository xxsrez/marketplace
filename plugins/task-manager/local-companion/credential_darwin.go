//go:build darwin

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

const keychainService = "com.xxsrez.task-manager.local-companion.oauth"

type keychainCredentialStore struct {
	account string
}

func (store *keychainCredentialStore) Load() (oauthCredential, error) {
	command := exec.Command(
		"/usr/bin/security", "find-generic-password", "-a", store.account,
		"-s", keychainService, "-w",
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		if strings.Contains(stderr.String(), "could not be found") ||
			strings.Contains(stderr.String(), "SecKeychainSearchCopyNext: The specified item could not be found") {
			return oauthCredential{}, errCredentialNotFound
		}
		return oauthCredential{}, newLocalError("keychain_unavailable", "OAuth credential could not be read from macOS Keychain")
	}
	var credential oauthCredential
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &credential); err != nil {
		return oauthCredential{}, newLocalError("keychain_invalid", "stored OAuth credential is invalid")
	}
	return credential, nil
}

func (store *keychainCredentialStore) Save(credential oauthCredential) error {
	secret, err := json.Marshal(credential)
	if err != nil {
		return newLocalError("keychain_invalid", "OAuth credential could not be encoded")
	}
	command := exec.Command(
		"/usr/bin/security", "add-generic-password", "-U", "-a", store.account,
		"-s", keychainService, "-l", "Task Manager local Codex companion", "-w",
	)
	// `security ... -w` prompts twice when no password argument is supplied.
	// Feed both confirmations through stdin so the secret never appears in argv.
	stdin := make([]byte, 0, len(secret)*2+2)
	stdin = append(stdin, secret...)
	stdin = append(stdin, '\n')
	stdin = append(stdin, secret...)
	stdin = append(stdin, '\n')
	command.Stdin = bytes.NewReader(stdin)
	command.Stdout = nil
	command.Stderr = nil
	if err := command.Run(); err != nil {
		return newLocalError("keychain_unavailable", "OAuth credential could not be stored in macOS Keychain")
	}
	return nil
}

func (store *keychainCredentialStore) Delete() error {
	command := exec.Command(
		"/usr/bin/security", "delete-generic-password", "-a", store.account,
		"-s", keychainService,
	)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		if strings.Contains(stderr.String(), "could not be found") {
			return nil
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return newLocalError("keychain_unavailable", "OAuth credential could not be removed from macOS Keychain")
		}
		return newLocalError("keychain_unavailable", "macOS Keychain command could not run")
	}
	return nil
}

type systemBrowser struct{}

func (systemBrowser) Open(rawURL string) error {
	return exec.Command("/usr/bin/open", rawURL).Run()
}
