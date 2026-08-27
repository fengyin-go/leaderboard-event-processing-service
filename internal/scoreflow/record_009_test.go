package scoreflow

import "testing"

func TestRecord009Lifecycle(t *testing.T) {
	got, err := RunScenario009()
	if err != nil || got != "confirmed-9" {
		t.Fatalf("record 009 lifecycle result mismatch: got %q err %v", got, err)
	}
}
