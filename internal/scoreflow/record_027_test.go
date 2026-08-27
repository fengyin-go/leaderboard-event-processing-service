package scoreflow

import "testing"

func TestRecord027Lifecycle(t *testing.T) {
	got, err := RunScenario027()
	if err != nil || got != "released" {
		t.Fatalf("record 027 lifecycle result mismatch: got %q err %v", got, err)
	}
}
