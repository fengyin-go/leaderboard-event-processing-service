package scoreflow

import "testing"

func TestRecord023Lifecycle(t *testing.T) {
	got, err := RunScenario023()
	if err != nil || got != "confirmed-23" {
		t.Fatalf("record 023 lifecycle result mismatch: got %q err %v", got, err)
	}
}
