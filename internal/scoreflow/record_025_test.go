package scoreflow

import "testing"

func TestRecord025Lifecycle(t *testing.T) {
	got, err := RunScenario025()
	if err != nil || got != "attempts=2,writes=1" {
		t.Fatalf("record 025 lifecycle result mismatch: got %q err %v", got, err)
	}
}
