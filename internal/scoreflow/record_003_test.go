package scoreflow

import "testing"

func TestRecord003Lifecycle(t *testing.T) {
	got, err := RunScenario003()
	if err != nil || got != "accepted-3" {
		t.Fatalf("record 003 lifecycle result mismatch: got %q err %v", got, err)
	}
}
