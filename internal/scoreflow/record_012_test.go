package scoreflow

import "testing"

func TestRecord012Lifecycle(t *testing.T) {
	got, err := RunScenario012()
	if err != nil || got != "missing" {
		t.Fatalf("record 012 lifecycle result mismatch: got %q err %v", got, err)
	}
}
