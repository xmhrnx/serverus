//go:build !windows

package monitor

import (
	"context"

	"github.com/shirou/gopsutil/v4/load"
)

// loadAverage is the 1, 5 and 15 minute load average.
func loadAverage(ctx context.Context) [3]float64 {
	avg, err := load.AvgWithContext(ctx)
	if err != nil {
		return [3]float64{}
	}
	return [3]float64{avg.Load1, avg.Load5, avg.Load15}
}
