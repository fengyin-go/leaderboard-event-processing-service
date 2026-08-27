package scoreflow

import "testing"

func TestRecord017Lifecycle(t *testing.T) {
	got, err := RunScenario017()
	if err != nil || got == "context canceled" {
		t.Fatalf("record 017 lifecycle result mismatch: got %q err %v", got, err)
	}
}
