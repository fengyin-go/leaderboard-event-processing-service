package scoreflow

import "testing"

func TestRecord002Lifecycle(t *testing.T) {
	got, err := RunScenario002()
	if err != nil || got != "confirmed-2" {
		t.Fatalf("record 002 lifecycle result mismatch: got %q err %v", got, err)
	}
}
