package scoreflow

import "testing"

func TestRecord006Lifecycle(t *testing.T) {
	got, err := RunScenario006()
	if err != nil || got != "released" {
		t.Fatalf("record 006 lifecycle result mismatch: got %q err %v", got, err)
	}
}
