//go:build !darwin && !linux

package main

import "os"

func openLocalFileNoFollow(string) (*os.File, error) {
	return nil, newLocalError("unsupported_platform", "local upload companion supports macOS and Linux only")
}

func platformFileIdentity(os.FileInfo) (stableFileIdentity, error) {
	return stableFileIdentity{}, newLocalError("unsupported_platform", "local upload companion supports macOS and Linux only")
}
