package overview

import (
	"fmt"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// temperatureTile is the warmest sensor.
func temperatureTile(hot monitor.Sensor) giode.Element {
	crit := hot.Critical
	if crit <= 0 {
		crit = 100
	}
	c := []giode.Color{theme.Success, theme.Warning, theme.Danger}[hot.Level()]
	return kpiTile("thermometer", "Temperature", format.Celsius(hot.Celsius), hot.Name,
		fmt.Sprintf("Critical at %s", format.Celsius(crit)), widgets.Ring(hot.Celsius/crit*100, 64, 7, c))
}
