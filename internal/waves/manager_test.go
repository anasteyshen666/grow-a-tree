package waves

import "testing"

func TestFirstWaveReleasesExpectedCount(t *testing.T) {
	m := NewManager()
	total := 0
	for i := 0; i < 25*60; i++ { // up to 25 simulated seconds at 60 TPS
		total += m.Update(1.0 / 60.0)
		if m.Number() == 1 && m.InPrep() { // wave 1 fully released, back to prep
			break
		}
	}
	if m.Number() != 1 {
		t.Fatalf("expected to reach wave 1, got %d", m.Number())
	}
	if total != waveSize(1) {
		t.Fatalf("wave 1 should release %d bugs, got %d", waveSize(1), total)
	}
}

func TestPrepDelaysFirstSpawn(t *testing.T) {
	m := NewManager()
	if n := m.Update(prepTime - 1); n != 0 || !m.InPrep() {
		t.Fatal("bugs spawned before prep ended")
	}
	if m.Number() != 0 {
		t.Fatalf("wave should still be 0 during prep, got %d", m.Number())
	}
}

func TestWaveSizeGrows(t *testing.T) {
	if waveSize(2) <= waveSize(1) {
		t.Fatal("later waves should be larger")
	}
}
