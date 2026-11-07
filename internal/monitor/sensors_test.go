package monitor

import (
	"testing"
)

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
