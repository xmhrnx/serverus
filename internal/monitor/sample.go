package monitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/sensors"
)

// Sample reads the local machine's counters.
func (HostSource) Sample(ctx context.Context) (Reading, error) {
	r := Reading{Time: time.Now()}

	if info, err := host.InfoWithContext(ctx); err == nil {
		r.System.Hostname = info.Hostname
		r.System.OS = strings.TrimSpace(fmt.Sprintf("%s %s", titleCase(info.Platform), info.PlatformVersion))
		r.System.Kernel = strings.TrimSpace(fmt.Sprintf("%s %s %s", titleCase(info.OS), info.KernelVersion, info.KernelArch))
		r.System.Uptime = time.Duration(info.Uptime) * time.Second
		r.Processes = int(info.Procs)
	}
	if infos, err := cpu.InfoWithContext(ctx); err == nil && len(infos) > 0 {
		r.System.CPUModel = strings.Join(strings.Fields(infos[0].ModelName), " ")
	}

	total, err := cpu.TimesWithContext(ctx, false)
	if err != nil {
		return Reading{}, err
	}
	if len(total) > 0 {
		r.CPU = fromCPUStat(total[0])
	}
	if cores, err := cpu.TimesWithContext(ctx, true); err == nil {
		for _, c := range cores {
			r.Cores = append(r.Cores, fromCPUStat(c))
		}
	}
	r.System.Cores = len(r.Cores)

	if vm, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		r.Memory = Memory{Total: vm.Total, Available: vm.Available, Free: vm.Free, Buffers: vm.Buffers, Cached: vm.Cached}
	}
	if sw, err := mem.SwapMemoryWithContext(ctx); err == nil {
		r.Memory.SwapTotal, r.Memory.SwapFree = sw.Total, sw.Free
	}
	r.Load = loadAverage(ctx)

	if parts, err := disk.PartitionsWithContext(ctx, false); err == nil {
		seen := map[string]bool{}
		for _, p := range parts {
			if ignoredFS[p.Fstype] || seen[p.Device] || strings.HasPrefix(p.Device, "/dev/loop") {
				continue
			}
			u, err := disk.UsageWithContext(ctx, p.Mountpoint)
			if err != nil || u.Total == 0 {
				continue
			}
			seen[p.Device] = true
			r.Disks = append(r.Disks, DiskUsage{Device: p.Device, FSType: p.Fstype, Mount: p.Mountpoint, Total: u.Total, Used: u.Used, Free: u.Free})
		}
	}
	if io, err := disk.IOCountersWithContext(ctx); err == nil {
		names := map[string]bool{}
		for n := range io {
			names[n] = true
		}
		for n, c := range io {
			if isPartition(n, names) || strings.HasPrefix(n, "loop") {
				continue
			}
			r.DiskIO = append(r.DiskIO, DiskCounter{Device: n, ReadBytes: c.ReadBytes, WriteBytes: c.WriteBytes})
		}
	}
	if ifs, err := net.IOCountersWithContext(ctx, true); err == nil {
		for _, c := range ifs {
			if c.Name == "lo" || strings.HasPrefix(c.Name, "Loopback") || c.Name == "lo0" {
				continue
			}
			r.Net = append(r.Net, NetCounter{Interface: c.Name, RxBytes: c.BytesRecv, TxBytes: c.BytesSent})
		}
	}
	if temps, err := sensors.TemperaturesWithContext(ctx); err == nil || len(temps) > 0 {
		for _, t := range temps {
			if t.Temperature <= 0 || t.Temperature >= 200 {
				continue
			}
			r.Sensors = append(r.Sensors, Sensor{Name: t.SensorKey, Celsius: t.Temperature, High: t.High, Critical: t.Critical})
		}
	}
	return r, nil
}
