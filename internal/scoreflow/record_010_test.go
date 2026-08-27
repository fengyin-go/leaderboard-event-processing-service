package scoreflow

import "testing"

func TestRecord010Lifecycle(t *testing.T) {
	got, err := RunScenario010()
	if err != nil || got != "accepted-10" {
		t.Fatalf("record 010 lifecycle result mismatch: got %q err %v", got, err)
	}
}
