//go:build !darwin

package main

import "os"

func openLocalFileNoFollow(string) (*os.File, error) {
	return nil, newLocalError("unsupported_platform", "local upload companion currently supports macOS only")
}

func platformFileIdentity(os.FileInfo) (stableFileIdentity, error) {
	return stableFileIdentity{}, newLocalError("unsupported_platform", "local upload companion currently supports macOS only")
}
