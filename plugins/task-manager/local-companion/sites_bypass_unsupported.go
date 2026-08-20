//go:build !darwin

package main

type keychainSitesBypassStore struct {
	account string
}

func (*keychainSitesBypassStore) LoadSecret() (string, error) {
	return "", newLocalError("unsupported_platform", "private UAT bridge currently supports macOS only")
}
