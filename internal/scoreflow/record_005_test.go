package scoreflow

import "testing"

func TestRecord005Lifecycle(t *testing.T) {
	got, err := RunScenario005()
	if err != nil || got != "missing" {
		t.Fatalf("record 005 lifecycle result mismatch: got %q err %v", got, err)
	}
}
