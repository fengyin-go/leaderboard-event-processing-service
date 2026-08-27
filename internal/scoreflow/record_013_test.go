package scoreflow

import "testing"

func TestRecord013Lifecycle(t *testing.T) {
	got, err := RunScenario013()
	if err != nil || got != "released" {
		t.Fatalf("record 013 lifecycle result mismatch: got %q err %v", got, err)
	}
}
