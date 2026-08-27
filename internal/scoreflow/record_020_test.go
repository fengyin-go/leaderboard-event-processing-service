package scoreflow

import "testing"

func TestRecord020Lifecycle(t *testing.T) {
	got, err := RunScenario020()
	if err != nil || got != "released" {
		t.Fatalf("record 020 lifecycle result mismatch: got %q err %v", got, err)
	}
}
