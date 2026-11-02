//go:build windows

package monitor

import "context"

// loadAverage is zero on Windows, which has no load average.
func loadAverage(context.Context) [3]float64 { return [3]float64{} }
