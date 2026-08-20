package main

import "os"

type stableFileIdentity struct {
	Device     uint64
	Inode      uint64
	Size       int64
	ModifiedNS int64
	ChangedNS  int64
}

func localFileIdentity(info os.FileInfo) (stableFileIdentity, error) {
	return platformFileIdentity(info)
}
