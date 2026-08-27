package scoreflow

import "testing"

func TestRecord001Lifecycle(t *testing.T) {
	got, err := RunScenario001()
	if err != nil || got != "confirmed-1" {
		t.Fatalf("record 001 lifecycle result mismatch: got %q err %v", got, err)
	}
}
