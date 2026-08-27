package scoreflow

import "testing"

func TestRecord018Lifecycle(t *testing.T) {
	got, err := RunScenario018()
	if err != nil || got != "attempts=2,writes=1" {
		t.Fatalf("record 018 lifecycle result mismatch: got %q err %v", got, err)
	}
}
