package ops

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// RootFiles changes files as root through sudo, for when the signed-in
// user isn't allowed to (files owned by root or another user).
type RootFiles struct {
	Exec transport.Executor
}

// IsPermissionError reports whether err is a file operation refused for
// lack of permission.
func IsPermissionError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, fs.ErrPermission) {
		return true
	}
	m := strings.ToLower(err.Error())
	return strings.Contains(m, "permission denied") || strings.Contains(m, "operation not permitted")
}

// Remove deletes path; a folder goes with everything in it.
func (r RootFiles) Remove(ctx context.Context, path string, isDir bool) error {
	flag := "-f"
	if isDir {
		flag = "-rf"
	}
	_, err := runAsRoot(ctx, r.Exec, "rm "+flag+" -- "+Quote(path))
	return err
}

// Rename moves from to to.
func (r RootFiles) Rename(ctx context.Context, from, to string) error {
	_, err := runAsRoot(ctx, r.Exec, "mv -- "+Quote(from)+" "+Quote(to))
	return err
}

// Chmod sets path's permission bits.
func (r RootFiles) Chmod(ctx context.Context, path string, mode fs.FileMode) error {
	_, err := runAsRoot(ctx, r.Exec, fmt.Sprintf("chmod %o -- %s", mode.Perm(), Quote(path)))
	return err
}
