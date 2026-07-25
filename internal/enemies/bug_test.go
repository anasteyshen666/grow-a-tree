package enemies

import "testing"

type mockField struct {
	adjacent bool
	onRot    bool
	rotStep  bool
	damaged  int
	rotEaten int
}

func (f *mockField) RandomEdgeSpawn() (int, int, bool)  { return 0, 0, true }
func (f *mockField) NextStep(c, r int) (int, int, bool) { return c + 1, r, true }
func (f *mockField) NextRotStep(c, r int) (int, int, bool) {
	if f.rotStep {
		return c, r + 1, true
	}
	return 0, 0, false
}
func (f *mockField) AdjacentTarget(c, r int) (int, int, bool) {
	if f.adjacent {
		return c + 1, r, true
	}
	return 0, 0, false
}
func (f *mockField) IsRot(int, int) bool { return f.onRot }
func (f *mockField) EatRot(int, int)     { f.rotEaten++ }
func (f *mockField) Damage(int, int)     { f.damaged++ }

func TestBugGnawsAdjacentTargetWithoutMoving(t *testing.T) {
	f := &mockField{adjacent: true}
	b := &Bug{Col: 5, Row: 5}
	b.update(gnawTime, f)

	if f.damaged != 1 {
		t.Fatalf("expected one gnaw, got %d", f.damaged)
	}
	if b.Col != 5 || b.Row != 5 {
		t.Fatal("bug moved while gnawing")
	}
}

func TestBugStepsWhenPathIsClear(t *testing.T) {
	f := &mockField{adjacent: false}
	b := &Bug{Col: 5, Row: 5}
	b.update(moveInterval, f)

	if b.Col != 6 || b.Row != 5 {
		t.Fatalf("bug did not step toward the network: (%d,%d)", b.Col, b.Row)
	}
	if f.damaged != 0 {
		t.Fatal("bug damaged something with no adjacent target")
	}
}

func TestBugChasesRotOverGnawing(t *testing.T) {
	f := &mockField{adjacent: true, rotStep: true} // root adjacent, but rot lures it
	b := &Bug{Col: 5, Row: 5}
	b.update(moveInterval, f)

	if b.Row != 6 || b.Col != 5 {
		t.Fatalf("bug did not head for the rot: (%d,%d)", b.Col, b.Row)
	}
	if f.damaged != 0 {
		t.Fatal("bug gnawed instead of chasing bait")
	}
}

func TestBugDiesOnRot(t *testing.T) {
	f := &mockField{onRot: true}
	b := &Bug{Col: 5, Row: 5}

	if dead := b.update(poisonTime, f); !dead {
		t.Fatal("bug should die after lingering on rot")
	}
	if f.rotEaten != 1 {
		t.Fatalf("bug did not consume the rot, eaten=%d", f.rotEaten)
	}
}
