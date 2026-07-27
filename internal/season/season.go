// Package season runs the meta-game cycle, tied to the wave number: a new
// season every wavesPerSeason waves, cycling summer -> autumn -> winter ->
// spring (GDD §6). Each season tweaks spores, water sources, rot, and the tint
// of the dirt background.
package season

type Season int

const (
	Summer Season = iota
	Autumn
	Winter
	Spring
)

const wavesPerSeason = 5

// Of returns the season for a wave number.
func Of(wave int) Season {
	if wave < 1 {
		return Summer
	}
	return Season(((wave - 1) / wavesPerSeason) % 4)
}

func (s Season) Name() string {
	switch s {
	case Autumn:
		return "AUTUMN"
	case Winter:
		return "WINTER"
	case Spring:
		return "SPRING"
	default:
		return "SUMMER"
	}
}

// SporeCap adjusts the spore cap: more in autumn, fewer in spring.
func (s Season) SporeCap(base int) int {
	switch s {
	case Autumn:
		return base + 2
	case Spring:
		if base-2 < 1 {
			return 1
		}
		return base - 2
	default:
		return base
	}
}

// SourceDrainMul scales how fast water sources deplete: faster in winter,
// slower in spring.
func (s Season) SourceDrainMul() float64 {
	switch s {
	case Winter:
		return 1.6
	case Spring:
		return 0.6
	default:
		return 1.0
	}
}

// RotLifetimeMul scales how long rot lasts: longer in winter.
func (s Season) RotLifetimeMul() float64 {
	if s == Winter {
		return 1.6
	}
	return 1.0
}

// DirtTint returns RGB multipliers applied to the dirt background: neutral in
// summer, warm orange in autumn, cold blue in winter, green in spring.
func (s Season) DirtTint() (r, g, b float64) {
	switch s {
	case Autumn:
		return 1.25, 0.95, 0.7
	case Winter:
		return 0.8, 0.95, 1.3
	case Spring:
		return 0.85, 1.15, 0.85
	default:
		return 1, 1, 1
	}
}
