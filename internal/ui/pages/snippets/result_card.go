package snippets

import (
	"fmt"
	"time"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// result is one server's outcome in a broadcast.
type result struct {
	serverID string
	name     string
	running  bool
	exitCode int
	output   string
	err      error
	took     time.Duration
}

// resultCard draws a server's outcome with its output.
func resultCard(r *result) giode.Element {
	var status giode.Element
	switch {
	case r.running:
		status = giode.Spinner(giode.Styles{Width: 16, Height: 16, Color: theme.Primary})
	case r.err != nil:
		status = widgets.Pill("error", theme.Danger)
	case r.exitCode != 0:
		status = widgets.Pill(fmt.Sprintf("exit %d", r.exitCode), theme.Warning)
	default:
		status = widgets.Pill("exit 0 · "+r.took.Round(time.Millisecond).String(), theme.Success)
	}
	body := r.output
	if r.err != nil {
		body = r.err.Error()
	}
	if body == "" && !r.running {
		body = "(no output)"
	}
	children := []giode.Element{
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(r.name, theme.Heading)), status),
	}
	if body != "" {
		children = append(children, giode.Box(giode.Styles{Width: -1, Background: theme.Sunken, BorderRadius: theme.RadiusSm, Padding: giode.UniformInset(theme.Space3), MaxHeight: 220, Overflow: giode.OverflowHidden},
			giode.Text(body, theme.Code.With(giode.Styles{LineHeight: 17})).Selectable("result-"+r.serverID)))
	}
	return widgets.Card(giode.Styles{Width: -1, Padding: giode.UniformInset(theme.Space4), Gap: theme.Space3}, children...)
}
