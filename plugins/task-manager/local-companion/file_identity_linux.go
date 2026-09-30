//go:build linux

package main

import (
	"os"
	"syscall"
)

func openLocalFileNoFollow(path string) (*os.File, error) {
	// Nonblocking open also prevents a FIFO replacement racing Lstat from hanging.
	return os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
}

func platformFileIdentity(info os.FileInfo) (stableFileIdentity, error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return stableFileIdentity{}, newLocalError(
			"local_path_unavailable", "authorized file metadata is unavailable",
		)
	}
	return stableFileIdentity{
		Device: uint64(stat.Dev), Inode: uint64(stat.Ino), Size: stat.Size,
		ModifiedNS: stat.Mtim.Sec*1_000_000_000 + stat.Mtim.Nsec,
		ChangedNS:  stat.Ctim.Sec*1_000_000_000 + stat.Ctim.Nsec,
	}, nil
}
