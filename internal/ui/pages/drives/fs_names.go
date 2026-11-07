package drives

// fsName is the file system's common name.
func fsName(t string) string {
	switch t {
	case "vfat":
		return "FAT32"
	case "exfat":
		return "exFAT"
	case "ntfs", "ntfs3":
		return "NTFS"
	case "crypto_LUKS":
		return "Encrypted"
	case "LVM2_member":
		return "LVM"
	case "swap":
		return "Swap"
	}
	return t
}
