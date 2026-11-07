package ops

import "testing"

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
