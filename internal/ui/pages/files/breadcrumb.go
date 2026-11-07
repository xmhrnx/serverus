package files

import (
	"strings"

	"github.com/xmhrnx/giode"
)

// Breadcrumb shows the current path as clickable segments.
type Breadcrumb struct {
	path     func() string
	navigate func(string)
	segments *giode.KeyedCache[string, *giode.PressableElement]
}

// NewBreadcrumb returns a breadcrumb for the path path() reports.
func NewBreadcrumb(path func() string, navigate func(string)) *Breadcrumb {
	return &Breadcrumb{path: path, navigate: navigate, segments: giode.NewKeyedCache[string, *giode.PressableElement]()}
}

// crumbs splits "/var/log/nginx" into its cumulative paths. Windows
// drive paths ("C:/Users") keep the drive as the first crumb.
func crumbs(p string) (labels, paths []string) {
	parts := strings.Split(strings.Trim(p, "/"), "/")
	prefix := ""
	if strings.HasPrefix(p, "/") {
		labels, paths, prefix = []string{"/"}, []string{"/"}, "/"
	}
	acc := strings.TrimSuffix(prefix, "/")
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 && prefix == "" {
			acc = part
		} else {
			acc += "/" + part
		}
		labels = append(labels, part)
		paths = append(paths, acc)
	}
	return labels, paths
}
