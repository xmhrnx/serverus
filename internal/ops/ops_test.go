package ops

import "testing"

func TestParsePS(t *testing.T) {
	out := "    1     0 root      0.1  0.2 12000  3600 Ss   systemd   /sbin/init splash\n" +
		"  842     1 www-data 12.5  1.4 98000   120 S    nginx     nginx: worker process\n"
	ps := parsePS(out)
	if len(ps) != 2 {
		t.Fatalf("got %d", len(ps))
	}
	if ps[1].PID != 842 || ps[1].User != "www-data" || ps[1].CPU != 12.5 || ps[1].Command != "nginx: worker process" {
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
