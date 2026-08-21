//go:build !darwin

package main

type systemBrowser struct{}

func (systemBrowser) Open(string) error {
	return newLocalError("unsupported_platform", "local upload companion currently supports macOS only")
}
