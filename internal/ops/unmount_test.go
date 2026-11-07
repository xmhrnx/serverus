package ops

import (
	"context"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// fakeMounts is a server where /dev/sdb1 is mounted on mounts, udisks
// is missing and umount needs root.
type fakeMounts struct {
	mounts []string
	ran    []string
}

// stuckMounts never lets go.
type stuckMounts struct{ fakeMounts }

func (f *fakeMounts) Run(_ context.Context, cmd string) (transport.Output, error) {
	f.ran = append(f.ran, cmd)
	switch {
	case strings.HasPrefix(cmd, "udisksctl"):
		return transport.Output{ExitCode: 127, Stderr: []byte("udisksctl: not found")}, nil
	case strings.HasPrefix(cmd, "sh -c") && strings.Contains(cmd, "/proc/mounts"):
		return transport.Output{Stdout: []byte(strings.Join(f.mounts, "\n"))}, nil
	case strings.HasPrefix(cmd, "umount"):
		return transport.Output{ExitCode: 32, Stderr: []byte("umount: must be superuser to use umount.")}, nil
	case strings.HasPrefix(cmd, "sudo -n sh -c 'umount "):
		if len(f.mounts) > 0 {
			f.mounts = f.mounts[:len(f.mounts)-1]
		}
		return transport.Output{}, nil
	}
	return transport.Output{ExitCode: 1, Stderr: []byte("unexpected " + cmd)}, nil
}
