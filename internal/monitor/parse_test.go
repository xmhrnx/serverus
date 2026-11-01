package monitor

import (
	"os"
	"testing"
	"time"
)

func TestParseProbeContainer(t *testing.T) {
	data, err := os.ReadFile("testdata/probe_container.txt")
	if err != nil {
		t.Fatal(err)
	}
	r, err := parseProbe(string(data), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if r.System.Hostname != "vm" || r.System.OS != "Ubuntu 24.04.5 LTS" {
		t.Errorf("system = %+v", r.System)
	}
	if r.System.Cores != 4 || len(r.Cores) != 4 {
		t.Errorf("cores = %d", r.System.Cores)
	}
	if r.Memory.Total != 16477352*1024 || r.Memory.Available != 15791696*1024 {
		t.Errorf("memory = %+v", r.Memory)
	}
	if len(r.Disks) != 1 || r.Disks[0].Mount != "/" {
		t.Errorf("disks = %+v (tmpfs must be skipped)", r.Disks)
	}
	if r.Processes != 80 {
		t.Errorf("processes = %d", r.Processes)
	}
	if r.Load[0] != 0.43 {
		t.Errorf("load = %v", r.Load)
	}
}

func TestParseProbeNotLinux(t *testing.T) {
	_, err := parseProbe("@@host\nmac\n@@os\n\nDarwin 23.1.0 arm64\n@@stat\n", time.Now())
	if err != ErrNotLinux {
		t.Fatalf("err = %v, want ErrNotLinux", err)
	}
}

func TestParseSensors(t *testing.T) {
	hw := []string{
		"coretemp|Package id 0|54000|80000|100000",
		"nvme|Composite|41850||84850",
		"acpitz||27800||",
		"bad|x|-273000||",
	}
	s := parseSensors(hw, []string{"x86_pkg_temp|60000"})
	if len(s) != 3 {
		t.Fatalf("got %d sensors: %+v", len(s), s)
	}
	if s[0].Name != "coretemp Package id 0" || s[0].Celsius != 54 || s[0].High != 80 || s[0].Critical != 100 {
		t.Errorf("sensor 0 = %+v", s[0])
	}
	if s[2].Name != "acpitz" {
		t.Errorf("unlabelled sensor = %+v", s[2])
	}
	if th := parseSensors(nil, []string{"x86_pkg_temp|60000"}); len(th) != 1 || th[0].Celsius != 60 {
		t.Errorf("thermal fallback = %+v", th)
	}
}

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

func TestSeriesRing(t *testing.T) {
	s := NewSeries(3)
	for i := 1; i <= 5; i++ {
		s.Push(float64(i))
	}
	got := s.Values()
	if len(got) != 3 || got[0] != 3 || got[2] != 5 {
		t.Fatalf("got %v", got)
	}
}
