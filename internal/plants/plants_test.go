package plants

import "testing"

func TestNoPlantsNeutralModifiers(t *testing.T) {
	m := NewManager()
	e, w := m.Modifiers()
	if e != 1 || w != 1 {
		t.Fatalf("empty manager should be neutral, got e=%v w=%v", e, w)
	}
}

func TestBatteryAndMossModifiers(t *testing.T) {
	m := NewManager()
	m.Add(0, 0, Battery)
	m.Add(1, 0, Moss)

	e, w := m.Modifiers()
	if e <= 1 {
		t.Fatalf("battery should raise energy multiplier, got %v", e)
	}
	if w >= 1 {
		t.Fatalf("moss should lower water multiplier, got %v", w)
	}
}

func TestWaterMultHasFloor(t *testing.T) {
	m := NewManager()
	for i := 0; i < 20; i++ { // pile on moss well past the floor
		m.Add(i, 0, Moss)
	}
	if _, w := m.Modifiers(); w < waterMultFloor {
		t.Fatalf("water multiplier dropped below the floor: %v", w)
	}
}
