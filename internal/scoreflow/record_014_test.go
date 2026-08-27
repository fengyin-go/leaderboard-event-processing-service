package scoreflow

import "testing"

func TestRecord014Lifecycle(t *testing.T) {
	got, err := RunScenario014()
	if err != nil || got != "first=true,second=false,changes=1" {
		t.Fatalf("record 014 lifecycle result mismatch: got %q err %v", got, err)
	}
}
