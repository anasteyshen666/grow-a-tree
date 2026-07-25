package resources

// Pool is a single resource balance clamped to [0, Max].
type Pool struct {
	Cur, Max float64
}

func (p *Pool) add(v float64) {
	p.Cur += v
	switch {
	case p.Cur > p.Max:
		p.Cur = p.Max
	case p.Cur < 0:
		p.Cur = 0
	}
}

const (
	RootEnergyCost = 8.0

	energyRegenPerSec = 6.0
	waterPerEnergy    = 0.5
	seedsRegenPerSec  = 2.0

	// Placeholder water inflow until real sources arrive (Этап 5).
	waterTricklePerSec = 2.0
)

type Resources struct {
	Energy Pool
	Water  Pool
	Seeds  Pool
}

func New() *Resources {
	return &Resources{
		Energy: Pool{Cur: 60, Max: 100},
		Water:  Pool{Cur: 100, Max: 100},
		Seeds:  Pool{Cur: 0, Max: 100},
	}
}

func (r *Resources) Update(dt float64) {
	r.Seeds.add(seedsRegenPerSec * dt)
	r.Water.add(waterTricklePerSec * dt)
	r.regenEnergy(dt)
}

// regenEnergy tops up Energy over time, but only while Water is available to
// fuel it — each point of Energy costs waterPerEnergy of Water (GDD §2).
func (r *Resources) regenEnergy(dt float64) {
	if r.Water.Cur <= 0 {
		return
	}
	gain := energyRegenPerSec * dt
	if space := r.Energy.Max - r.Energy.Cur; gain > space {
		gain = space
	}
	if cost := gain * waterPerEnergy; cost > r.Water.Cur {
		gain = r.Water.Cur / waterPerEnergy
	}
	r.Energy.add(gain)
	r.Water.add(-gain * waterPerEnergy)
}

// TrySpendEnergy deducts v if affordable and reports success.
func (r *Resources) TrySpendEnergy(v float64) bool {
	if r.Energy.Cur < v {
		return false
	}
	r.Energy.add(-v)
	return true
}
