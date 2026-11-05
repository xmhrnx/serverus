package ops

import "testing"

func TestParsePS(t *testing.T) {
	out := "    1     0 root      0.1  0.2 12000  3600 Ss   ?     systemd   /sbin/init splash\n" +
		"  842     1 www-data 12.5  1.4 98000   120 S    pts/0 nginx     nginx: worker process\n"
	ps := parsePS(out)
	if len(ps) != 2 {
		t.Fatalf("got %d", len(ps))
	}
	if ps[1].PID != 842 || ps[1].User != "www-data" || ps[1].CPU != 12.5 || ps[1].Command != "nginx: worker process" || ps[1].TTY != "pts/0" || ps[0].TTY != "" {
		t.Errorf("got %+v", ps[1])
	}
	if ps[0].RSS != 12000*1024 || ps[0].Elapsed.Hours() != 1 {
		t.Errorf("got %+v", ps[0])
	}
}

func TestParseServices(t *testing.T) {
	units := "nginx.service loaded active running A high performance web server\n" +
		"bad.service loaded failed failed Broken thing\n"
	files := "nginx.service enabled enabled\nbad.service disabled enabled\n"
	s := parseServices(units, files)
	if len(s) != 2 || !s[0].Running() || s[0].Enabled != "enabled" || !s[1].Failed() {
		t.Fatalf("got %+v", s)
	}
	if s[0].Description != "A high performance web server" {
		t.Errorf("description = %q", s[0].Description)
	}
}

func TestParseContainers(t *testing.T) {
	list := `{"ID":"abc","Names":"web","Image":"nginx:1.27","State":"running","Status":"Up 2 hours","Ports":"0.0.0.0:80->80/tcp"}
{"ID":"def","Names":"job","Image":"alpine","Status":"Exited (0) 3 days ago"}`
	stats := `{"ID":"abc","Name":"web","CPUPerc":"0.15%","MemUsage":"12MiB / 1GiB"}`
	cs := parseContainers(list, stats)
	if len(cs) != 2 || !cs[0].Running() || cs[0].CPU != "0.15%" || cs[1].State != "exited" {
		t.Fatalf("got %+v", cs)
	}
}

func TestQuote(t *testing.T) {
	if got := Quote("it's"); got != `'it'\''s'` {
		t.Errorf("got %s", got)
	}
}

func TestWithoutOwnCommands(t *testing.T) {
	procs := []Process{
		{PID: 900, PPID: 1, Name: "sshd"},
		{PID: 1000, PPID: 900, Name: "sshd-session"},        // this app's connection
		{PID: 1001, PPID: 1000, Name: "bash", TTY: "pts/1"}, // a terminal tab
		{PID: 1002, PPID: 1001, Name: "htop", TTY: "pts/1"},
		{PID: 1010, PPID: 1000, Name: "sh"},   // the listing
		{PID: 1011, PPID: 1010, Name: "ps"},   // its ps
		{PID: 1020, PPID: 1000, Name: "bash"}, // the metrics probe
		{PID: 1021, PPID: 1020, Name: "cat"},
		{PID: 2000, PPID: 900, Name: "sshd-session"}, // someone else's connection
		{PID: 2001, PPID: 2000, Name: "rsync"},
		{PID: 3000, PPID: 1, Name: "nginx"},
	}
	got := withoutOwnCommands(procs, 1010)
	names := []int32{}
	for _, p := range got {
		names = append(names, p.PID)
	}
	want := []int32{900, 1000, 1001, 1002, 2000, 2001, 3000}
	if len(names) != len(want) {
		t.Fatalf("kept %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("kept %v, want %v", names, want)
		}
	}
}
