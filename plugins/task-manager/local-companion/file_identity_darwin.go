//go:build darwin

package main

import (
	"os"
	"syscall"
)

func openLocalFileNoFollow(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
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
		ModifiedNS: stat.Mtimespec.Sec*1_000_000_000 + stat.Mtimespec.Nsec,
		ChangedNS:  stat.Ctimespec.Sec*1_000_000_000 + stat.Ctimespec.Nsec,
	}, nil
}
