package ops

import "testing"

func TestUnitFromCgroup(t *testing.T) {
	for in, want := range map[string]string{
		"0::/system.slice/nginx.service\n":                                     "nginx.service",
		"0::/user.slice/user-1000.slice/session-3.scope\n":                     "session-3.scope",
		"12:pids:/\n1:name=systemd:/system.slice/docker.service\n0::/\n":       "docker.service",
		"0::/user.slice/user-1000.slice/user@1000.service/app.slice/x.service": "user@1000.service",
		"0::/\n": "",
	} {
		if got := unitFromCgroup(in); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
	if (Process{Unit: "user@1000.service"}).Service() != "" {
		t.Error("the user manager isn't a service to stop")
	}
}
