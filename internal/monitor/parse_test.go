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
