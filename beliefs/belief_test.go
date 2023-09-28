package beliefs

import (
	"testing"
)

func TestNew(t *testing.T) {
	b1 := New("b1")
	b2 := New("b2")

	if b1.Name != "b1" {
		t.Errorf("b1.Name is not b1 it is %s", b1.Name)
	}

	if b2.Name != "b2" {
		t.Errorf("b2.Name is not b2 it is %s", b2.Name)
	}

	if b1.Uuid == b2.Uuid {
		t.Errorf("UUIDs equal when they should not be. b1.Uuid=%v; b2.Uuid=%v", b1.Uuid, b2.Uuid)
	}
}
