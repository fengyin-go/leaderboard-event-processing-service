package scoreflow

import "testing"

func TestRecord022Lifecycle(t *testing.T) {
	got, err := RunScenario022()
	if err != nil || got != "confirmed-22" {
		t.Fatalf("record 022 lifecycle result mismatch: got %q err %v", got, err)
	}
}
