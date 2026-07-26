package resources

import "testing"

func TestTrySpendEnergy(t *testing.T) {
	r := New()
	r.Energy.Cur = 10
	if !r.TrySpendEnergy(8) || r.Energy.Cur != 2 {
		t.Fatalf("spend failed, energy=%v", r.Energy.Cur)
	}
	if r.TrySpendEnergy(8) {
		t.Fatal("spent energy it did not have")
	}
}

func TestEnergyRegenConsumesWater(t *testing.T) {
	r := New()
	r.Energy.Cur, r.Water.Cur = 0, 100
	before := r.Water.Cur
	r.regenEnergy(1, 1, 1)
	if r.Energy.Cur <= 0 {
		t.Fatal("energy did not regenerate")
	}
	if r.Water.Cur >= before {
		t.Fatal("regen did not consume water")
	}
}

func TestNoEnergyRegenWithoutWater(t *testing.T) {
	r := New()
	r.Energy.Cur, r.Water.Cur = 0, 0
	r.regenEnergy(1, 1, 1)
	if r.Energy.Cur != 0 {
		t.Fatalf("energy regenerated without water: %v", r.Energy.Cur)
	}
}

func TestBatteryAuraSpeedsRegen(t *testing.T) {
	base, boosted := New(), New()
	base.Energy.Cur, boosted.Energy.Cur = 0, 0
	base.regenEnergy(1, 1.0, 1)
	boosted.regenEnergy(1, 1.2, 1) // +20% battery aura
	if boosted.Energy.Cur <= base.Energy.Cur {
		t.Fatal("battery aura did not speed up energy regen")
	}
}

func TestTrySpendSeeds(t *testing.T) {
	r := New()
	r.Seeds.Cur = 30
	if !r.TrySpendSeeds(20) || r.Seeds.Cur != 10 {
		t.Fatalf("seed spend failed, seeds=%v", r.Seeds.Cur)
	}
	if r.TrySpendSeeds(20) {
		t.Fatal("spent seeds it did not have")
	}
}

func TestPoolClamps(t *testing.T) {
	p := Pool{Cur: 95, Max: 100}
	p.add(20)
	if p.Cur != 100 {
		t.Fatalf("pool exceeded max: %v", p.Cur)
	}
	p.add(-200)
	if p.Cur != 0 {
		t.Fatalf("pool went below zero: %v", p.Cur)
	}
}
