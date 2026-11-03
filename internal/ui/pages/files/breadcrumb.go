package files

import (
	"strings"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
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

// Layout draws the segments.
func (b *Breadcrumb) Layout(gtx layout.Context) layout.Dimensions {
	labels, paths := crumbs(b.path())
	items := []giode.Element{}
	for i := range labels {
		label, target := labels[i], paths[i]
		last := i == len(labels)-1
		seg := b.segments.Get(target, func() *giode.PressableElement {
			return giode.Pressable(func(ps giode.PressState) giode.Element {
				fg, bg := theme.Muted, giode.Color("transparent")
				if ps.Hovered {
					fg, bg = theme.Text, theme.HoverFill
				}
				return giode.Box(giode.Styles{Padding: giode.SymmetricInset(6, 3), BorderRadius: theme.RadiusSm, Background: bg},
					giode.Text(label, giode.Styles{Color: fg, FontSize: theme.TextSm, FontFamily: theme.Mono}))
			}).Label(label).OnClick(func(giode.ClickEvent) { b.navigate(target) })
		})
		if last {
			items = append(items, giode.Box(giode.Styles{Padding: giode.SymmetricInset(6, 3)},
				giode.Text(label, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontFamily: theme.Mono, FontWeight: giode.FontWeightSemiBold})))
		} else {
			items = append(items, seg)
		}
		if !last && label != "/" {
			items = append(items, theme.Icon("chevron-right", 13, theme.Faint))
		}
	}
	return giode.Row(giode.Styles{Gap: 0, Align: giode.AlignCenter}, items...).Layout(gtx)
}
