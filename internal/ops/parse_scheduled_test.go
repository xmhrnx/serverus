package ops

import (
	"testing"
	"time"
)

func TestParseScheduled(t *testing.T) {
	s, ok := parseScheduled("USEC=1762279200000000\nWARN_WALL=1\nMODE=reboot\n")
	if !ok || s.Action != PowerReboot || !s.At.Equal(time.UnixMicro(1762279200000000)) {
		t.Fatalf("got %+v %v", s, ok)
	}
	if _, ok := parseScheduled(""); ok {
		t.Error("nothing scheduled")
	}
}
