package ops

import (
	"context"
	"fmt"
)

// Mkdir creates the folder path.
func (r RootFiles) Mkdir(ctx context.Context, path string) error {
	_, err := runAsRoot(ctx, r.Exec, "mkdir -- "+Quote(path))
	return err
}

// Write replaces path's contents with data, creating it if needed. An
// existing file keeps its owner and permissions.
func (r RootFiles) Write(ctx context.Context, path string, data []byte) error {
	_, err := runAsRootInput(ctx, r.Exec, "cat > "+Quote(path), data)
	return err
}

// Read returns path's contents, up to limit bytes.
func (r RootFiles) Read(ctx context.Context, path string, limit int) ([]byte, error) {
	out, err := runAsRoot(ctx, r.Exec, fmt.Sprintf("head -c %d -- %s", limit, Quote(path)))
	return []byte(out), err
}
