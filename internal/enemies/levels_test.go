package enemies

import "testing"

// maxLevelSeen samples levelForWave many times and returns the highest level it
// ever produces for the given wave.
func maxLevelSeen(wave int) int {
	top := 0
	for i := 0; i < 500; i++ {
		if l := levelForWave(wave); l > top {
			top = l
		}
	}
	return top
}

func TestLevelUnlockEveryTwoWavesWithReset(t *testing.T) {
	want := map[int]int{
		1: 1, 2: 1, // level 1 only
		3: 2, 4: 2, // + level 2
		5: 3, 6: 3,
		7: 4, 8: 4,
		9: 5, 10: 5, // full range
		11: 1, 12: 1, // resets to level 1
		13: 2, // cycle repeats
	}
	for wave, wantTop := range want {
		if got := maxLevelSeen(wave); got != wantTop {
			t.Errorf("wave %d: top level %d, want %d", wave, got, wantTop)
		}
	}
}

func TestLevelsAreMixed(t *testing.T) {
	// By wave 9 the full 1..5 range is unlocked; every level should appear.
	seen := map[int]bool{}
	for i := 0; i < 2000; i++ {
		seen[levelForWave(9)] = true
	}
	for lvl := 1; lvl <= maxLevel; lvl++ {
		if !seen[lvl] {
			t.Fatalf("level %d never spawned at wave 9 (not mixed)", lvl)
		}
	}
}
