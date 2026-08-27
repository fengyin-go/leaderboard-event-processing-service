package scoreflow

import "testing"

func TestRecord007Lifecycle(t *testing.T) {
	got, err := RunScenario007()
	if err != nil || got != "first=true,second=false,changes=1" {
		t.Fatalf("record 007 lifecycle result mismatch: got %q err %v", got, err)
	}
}
