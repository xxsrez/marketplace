//go:build !darwin

package main

type keychainCredentialStore struct{ account string }

func (*keychainCredentialStore) Load() (oauthCredential, error) {
	return oauthCredential{}, newLocalError("unsupported_platform", "local upload companion currently supports macOS only")
}
func (*keychainCredentialStore) Save(oauthCredential) error {
	return newLocalError("unsupported_platform", "local upload companion currently supports macOS only")
}
func (*keychainCredentialStore) Delete() error { return nil }

type systemBrowser struct{}

func (systemBrowser) Open(string) error {
	return newLocalError("unsupported_platform", "local upload companion currently supports macOS only")
}
