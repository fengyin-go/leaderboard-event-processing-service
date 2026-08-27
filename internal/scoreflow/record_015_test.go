package scoreflow

import "testing"

func TestRecord015Lifecycle(t *testing.T) {
	got, err := RunScenario015()
	if err != nil || got != "confirmed-15" {
		t.Fatalf("record 015 lifecycle result mismatch: got %q err %v", got, err)
	}
}
