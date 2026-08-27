package scoreflow

import "testing"

func TestRecord030Lifecycle(t *testing.T) {
	got, err := RunScenario030()
	if err != nil || got != "confirmed-30" {
		t.Fatalf("record 030 lifecycle result mismatch: got %q err %v", got, err)
	}
}
