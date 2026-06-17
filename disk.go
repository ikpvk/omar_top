package main

import (
	"bufio"
	"os"
	"strings"
	"syscall"
)

func readDisks() []DiskMetrics {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return nil
	}
	defer f.Close()

	var disks []DiskMetrics
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		device := fields[0]
		mountpoint := fields[1]
		fstype := fields[2]

		if !isPhysicalFS(fstype) {
			continue
		}
		if !strings.HasPrefix(device, "/") {
			continue
		}
		if isVirtualMount(mountpoint) {
			continue
		}

		var stat syscall.Statfs_t
		if err := syscall.Statfs(mountpoint, &stat); err != nil {
			continue
		}

		total := stat.Blocks * uint64(stat.Bsize)
		free := stat.Bfree * uint64(stat.Bsize)
		used := total - free

		disks = append(disks, DiskMetrics{
			Mountpoint: mountpoint,
			Fstype:     fstype,
			UsedBytes:  used,
			TotalBytes: total,
		})
	}

	return disks
}

func isPhysicalFS(fstype string) bool {
	switch fstype {
	case "ext4", "ext3", "ext2", "btrfs", "xfs", "zfs", "ntfs", "vfat", "fuseblk":
		return true
	}
	return false
}

func isVirtualMount(mountpoint string) bool {
	switch mountpoint {
	case "/proc", "/sys", "/dev", "/run", "/tmp":
		return true
	}
	return strings.HasPrefix(mountpoint, "/var") ||
		strings.HasPrefix(mountpoint, "/sys/") ||
		strings.HasPrefix(mountpoint, "/proc/") ||
		strings.HasPrefix(mountpoint, "/dev/") ||
		strings.HasPrefix(mountpoint, "/run/")
}
