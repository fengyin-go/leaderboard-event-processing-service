package scoreflow

import "testing"

func TestRecord004Lifecycle(t *testing.T) {
	got, err := RunScenario004()
	if err != nil || got != "attempts=2,writes=1" {
		t.Fatalf("record 004 lifecycle result mismatch: got %q err %v", got, err)
	}
}
