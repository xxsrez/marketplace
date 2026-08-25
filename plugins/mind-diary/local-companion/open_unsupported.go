//go:build !darwin

package main

import (
	"errors"
	"os"
)

func openRegularFileNoFollow(_ string) (*os.File, error) {
	return nil, errors.New("unsupported platform")
}
