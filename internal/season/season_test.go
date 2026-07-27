package season

import "testing"

func TestSeasonCycle(t *testing.T) {
	cases := map[int]Season{
		0: Summer, 1: Summer, 5: Summer,
		6: Autumn, 10: Autumn,
		11: Winter, 15: Winter,
		16: Spring, 20: Spring,
		21: Summer, // wraps around
	}
	for wave, want := range cases {
		if got := Of(wave); got != want {
			t.Errorf("wave %d: got %v, want %v", wave, got, want)
		}
	}
}

func TestSporeCap(t *testing.T) {
	if Autumn.SporeCap(4) <= 4 {
		t.Fatal("autumn should raise the spore cap")
	}
	if Spring.SporeCap(4) >= 4 {
		t.Fatal("spring should lower the spore cap")
	}
	if Summer.SporeCap(4) != 4 {
		t.Fatal("summer should not change the spore cap")
	}
}

func TestSourceDrainAndRot(t *testing.T) {
	if Winter.SourceDrainMul() <= 1 || Spring.SourceDrainMul() >= 1 {
		t.Fatal("winter should drain faster, spring slower")
	}
	if Winter.RotLifetimeMul() <= 1 {
		t.Fatal("winter rot should last longer")
	}
	if Summer.RotLifetimeMul() != 1 {
		t.Fatal("summer rot lifetime should be unchanged")
	}
}
