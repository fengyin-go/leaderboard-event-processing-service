package scoreflow

import "testing"

func TestRecord008Lifecycle(t *testing.T) {
	got, err := RunScenario008()
	if err != nil || got != "confirmed-8" {
		t.Fatalf("record 008 lifecycle result mismatch: got %q err %v", got, err)
	}
}
