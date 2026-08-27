package scoreflow

import "testing"

func TestRecord024Lifecycle(t *testing.T) {
	got, err := RunScenario024()
	if err != nil || got == "context canceled" {
		t.Fatalf("record 024 lifecycle result mismatch: got %q err %v", got, err)
	}
}
