package scoreflow

import "testing"

func TestRecord028Lifecycle(t *testing.T) {
	got, err := RunScenario028()
	if err != nil || got != "first=true,second=false,changes=1" {
		t.Fatalf("record 028 lifecycle result mismatch: got %q err %v", got, err)
	}
}
