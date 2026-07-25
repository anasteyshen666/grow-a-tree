package enemies

import "testing"

type mockField struct {
	adjacent bool
	damaged  int
}

func (f *mockField) RandomEdgeSpawn() (int, int, bool)     { return 0, 0, true }
func (f *mockField) NextStep(c, r int) (int, int, bool)    { return c + 1, r, true }
func (f *mockField) AdjacentTarget(c, r int) (int, int, bool) {
	if f.adjacent {
		return c + 1, r, true
	}
	return 0, 0, false
}
func (f *mockField) Damage(int, int) { f.damaged++ }

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
