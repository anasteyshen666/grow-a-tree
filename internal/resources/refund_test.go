package resources

import (
	"math"
	"testing"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func TestRefundRisesEvery10Waves(t *testing.T) {
	if !approx(RefundFor(0), RootEnergyCost*0.2) {
		t.Fatalf("wave 0 refund should be 20%%, got %v", RefundFor(0))
	}
	if !approx(RefundFor(9), RootEnergyCost*0.2) {
		t.Fatal("refund should not rise until wave 10")
	}
	if !approx(RefundFor(10), RootEnergyCost*0.3) {
		t.Fatalf("wave 10 refund should be 30%%, got %v", RefundFor(10))
	}
	if RefundFor(1000) > RootEnergyCost*0.9+1e-6 {
		t.Fatal("refund should cap at 90%")
	}
}
