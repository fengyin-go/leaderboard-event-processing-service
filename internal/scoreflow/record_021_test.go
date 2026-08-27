package scoreflow

import "testing"

func TestRecord021Lifecycle(t *testing.T) {
	got, err := RunScenario021()
	if err != nil || got != "first=true,second=false,changes=1" {
		t.Fatalf("record 021 lifecycle result mismatch: got %q err %v", got, err)
	}
}
