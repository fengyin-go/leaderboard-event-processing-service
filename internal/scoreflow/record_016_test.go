package scoreflow

import "testing"

func TestRecord016Lifecycle(t *testing.T) {
	got, err := RunScenario016()
	if err != nil || got != "confirmed-16" {
		t.Fatalf("record 016 lifecycle result mismatch: got %q err %v", got, err)
	}
}
