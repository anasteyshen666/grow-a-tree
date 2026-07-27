package waves

import "testing"

func TestPrepDelaysFirstSpawn(t *testing.T) {
	m := NewManager()
	if n := m.Update(prepTime-1, 0); n != 0 || !m.InPrep() {
		t.Fatal("bugs spawned before prep ended")
	}
	if m.Number() != 0 {
		t.Fatalf("wave should still be 0 during prep, got %d", m.Number())
	}
}

func TestFirstWaveReleasesExpectedCount(t *testing.T) {
	m := NewManager()
	m.Update(prepTime, 0)      // finish prep -> wave 1 starts spawning
	total := m.Update(1000, 0) // large dt drains the whole wave at once
	if m.Number() != 1 {
		t.Fatalf("expected wave 1, got %d", m.Number())
	}
	if total != waveSize(1) {
		t.Fatalf("wave 1 should release %d bugs, got %d", waveSize(1), total)
	}
}

func TestNextWaveWaitsForBugsCleared(t *testing.T) {
	m := NewManager()
	m.Update(prepTime, 0) // -> spawning wave 1
	m.Update(1000, 5)     // release the whole wave; now clearing, 5 bugs alive

	m.Update(1000, 3) // bugs still alive: prep must not start, however long we wait
	if m.InPrep() {
		t.Fatal("prep started while bugs from the wave were still alive")
	}

	m.Update(0, 0) // field cleared
	if !m.InPrep() || m.PrepRemaining() != prepTime {
		t.Fatal("prep did not start after the field was cleared")
	}
}

func TestWaveSizeGrows(t *testing.T) {
	if waveSize(1) != waveSize(2) {
		t.Fatal("count should be constant within a 2-wave step")
	}
	if waveSize(3) <= waveSize(1) {
		t.Fatal("count should step up every 2 waves")
	}
	// same step (waves 1-2 vs 11-12), later decade should have more bugs
	if waveSize(11) <= waveSize(1) {
		t.Fatal("each decade should have more bugs at the same step")
	}
	// each decade should also grow faster than the last
	d0 := waveSize(3) - waveSize(1)
	d1 := waveSize(13) - waveSize(11)
	if d1 <= d0 {
		t.Fatal("later decades should ramp up faster")
	}
}
