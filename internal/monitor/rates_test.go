package monitor

import (
	"testing"
	"time"
)

func TestParseDiskStatsSkipsPartitions(t *testing.T) {
	lines := []string{
		"   8       0 sda 100 0 2000 0 50 0 4000 0 0 0 0",
		"   8       1 sda1 90 0 1800 0 40 0 3000 0 0 0 0",
		" 259       0 nvme0n1 10 0 200 0 5 0 400 0 0 0 0",
		" 259       1 nvme0n1p1 10 0 200 0 5 0 400 0 0 0 0",
		"   7       0 loop0 1 0 2 0 0 0 0 0 0 0 0",
	}
	got := parseDiskStats(lines)
	if len(got) != 2 || got[0].Device != "sda" || got[1].Device != "nvme0n1" {
		t.Fatalf("got %+v", got)
	}
	if got[0].ReadBytes != 2000*512 || got[0].WriteBytes != 4000*512 {
		t.Errorf("sda = %+v", got[0])
	}
}

func TestDeriveRates(t *testing.T) {
	t0 := time.Unix(0, 0)
	prev := Reading{Time: t0, CPU: CPUTimes{User: 100, Idle: 900},
		Net: []NetCounter{{Interface: "eth0", RxBytes: 1000, TxBytes: 500}}}
	cur := Reading{Time: t0.Add(2 * time.Second), CPU: CPUTimes{User: 150, Idle: 950},
		Net: []NetCounter{{Interface: "eth0", RxBytes: 3000, TxBytes: 700}}}
	s := derive(cur, &prev)
	if s.CPUPercent != 50 {
		t.Errorf("cpu = %v, want 50", s.CPUPercent)
	}
	if s.NetRxRate != 1000 || s.NetTxRate != 100 {
		t.Errorf("net = %v / %v", s.NetRxRate, s.NetTxRate)
	}
}
