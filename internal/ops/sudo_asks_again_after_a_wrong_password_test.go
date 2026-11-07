package ops

import (
	"context"
	"testing"
)

func TestSudoAsksAgainAfterAWrongPassword(t *testing.T) {
	f := &fakeSudo{answers: []string{"nope", "pw"}}
	out, err := runPrivileged(context.Background(), f, "shutdown -c")
	if err != nil || out != "ok" {
		t.Fatalf("got %q, %v", out, err)
	}
	if len(f.asked) != 2 || f.asked[0] || !f.asked[1] {
		t.Errorf("asked %v", f.asked)
	}
	if last := f.commands[len(f.commands)-1]; last != "sudo -S -p '' sh -c 'shutdown -c'" {
		t.Errorf("ran %q", last)
	}
}

func TestSudoDeclined(t *testing.T) {
	f := &fakeSudo{}
	if _, err := runPrivileged(context.Background(), f, "shutdown -c"); err != ErrCanceled {
		t.Fatalf("got %v", err)
	}
}

func TestCheckMountDir(t *testing.T) {
	for dir, ok := range map[string]bool{
		"/mnt/data": true, "/srv/backup": true, "/home/pi/usb": true,
		"/": false, "/etc": false, "/mnt": false, "/home/": false, "data": false, "": false,
	} {
		if err := CheckMountDir(dir); (err == nil) != ok {
			t.Errorf("%q: %v", dir, err)
		}
	}
}
