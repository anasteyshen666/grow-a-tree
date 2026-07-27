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

	// Energy regenerates slowly; each point still drains the same water per
	// second as before (regen*waterPerEnergy is unchanged at 3/sec).
	energyRegenPerSec = 3.0
	waterPerEnergy    = 1.0
	seedsRegenPerSec  = 1.0
)

// RefundFor is the energy returned when cutting a root: 20% at the start, rising
// 10% every 10 waves (capped at 90%), rewarding late-game rot play.
func RefundFor(wave int) float64 {
	frac := 0.2 + 0.1*float64(wave/10)
	if frac > 0.9 {
		frac = 0.9
	}
	return RootEnergyCost * frac
}

type Resources struct {
	Energy Pool
	Water  Pool
	Seeds  Pool
}

const (
	BaseWaterMax  = 100.0
	BaseEnergyMax = 100.0
	PerCoreLink   = 50.0 // each core-to-core link raises the water and energy caps

	// WaterUpkeepPerCore is water/sec drained per extra Core: more trees give
	// more storage but cost more to sustain.
	WaterUpkeepPerCore = 3.0
)

func New() *Resources {
	return &Resources{
		Energy: Pool{Cur: BaseEnergyMax, Max: BaseEnergyMax},
		Water:  Pool{Cur: BaseWaterMax, Max: BaseWaterMax},
		Seeds:  Pool{Cur: 0, Max: 100},
	}
}

// SetCoreLinks raises the water and energy caps by the number of core-to-core
// links (GDD §2).
func (r *Resources) SetCoreLinks(merges int) {
	bonus := float64(merges) * PerCoreLink
	r.Water.Max = BaseWaterMax + bonus
	r.Energy.Max = BaseEnergyMax + bonus
}

// Update ticks the economy. energyMult and waterMult are plant-aura modifiers
// (1 = no plants): energyMult speeds regen, waterMult scales its water cost.
// waterUpkeep is water/sec drained passively by extra Cores.
func (r *Resources) Update(dt, energyMult, waterMult, waterUpkeep float64) {
	r.Seeds.add(seedsRegenPerSec * dt)
	if waterUpkeep > 0 {
		r.Water.add(-waterUpkeep * dt)
	}
	r.regenEnergy(dt, energyMult, waterMult)
}

// regenEnergy tops up Energy over time, but only while Water is available to
// fuel it — each point of Energy costs waterPerEnergy of Water (GDD §2).
func (r *Resources) regenEnergy(dt, energyMult, waterMult float64) {
	if r.Water.Cur <= 0 {
		return
	}
	costPer := waterPerEnergy * waterMult
	gain := energyRegenPerSec * energyMult * dt
	if space := r.Energy.Max - r.Energy.Cur; gain > space {
		gain = space
	}
	if cost := gain * costPer; cost > r.Water.Cur {
		gain = r.Water.Cur / costPer
	}
	r.Energy.add(gain)
	r.Water.add(-gain * costPer)
}

// AddEnergy credits energy (clamped to Max), e.g. a refund from cutting a root.
func (r *Resources) AddEnergy(v float64) {
	r.Energy.add(v)
}

// AddWater credits mined water (clamped to Max).
func (r *Resources) AddWater(v float64) {
	r.Water.add(v)
}

// TrySpendEnergy deducts v if affordable and reports success.
func (r *Resources) TrySpendEnergy(v float64) bool {
	if r.Energy.Cur < v {
		return false
	}
	r.Energy.add(-v)
	return true
}

// TrySpendSeeds deducts v seeds if affordable and reports success.
func (r *Resources) TrySpendSeeds(v float64) bool {
	if r.Seeds.Cur < v {
		return false
	}
	r.Seeds.add(-v)
	return true
}
