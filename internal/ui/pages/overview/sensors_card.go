package overview

import (
	"fmt"
	"sort"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// maxSensorsShown keeps the card a reasonable height on servers with
// dozens of probes.
const maxSensorsShown = 12

// SensorsCard shows every temperature, warmest first, in columns across
// the card. The dashboard only shows it when the machine has sensors.
func SensorsCard(s monitor.Snapshot) giode.Element {
	sensors := append([]monitor.Sensor(nil), s.Sensors...)
	sort.Slice(sensors, func(i, j int) bool { return sensors[i].Celsius > sensors[j].Celsius })
	more := ""
	if len(sensors) > maxSensorsShown {
		more = fmt.Sprintf("+ %d more", len(sensors)-maxSensorsShown)
		sensors = sensors[:maxSensorsShown]
	}
	cells := make([]giode.Element, len(sensors))
	for i, x := range sensors {
		cells[i] = sensorReading(x)
	}
	children := []giode.Element{
		widgets.CardHeader("thermometer", theme.Warning, "Temperatures", fmt.Sprintf("%d sensors", len(s.Sensors))),
		widgets.Grid(200, theme.Space5, cells...),
	}
	if more != "" {
		children = append(children, giode.Text(more, theme.Caption))
	}
	return widgets.Card(giode.Styles{Width: -1}, children...)
}
