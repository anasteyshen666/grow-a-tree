package world

// FxKind is a visual event the grid reports for the effects layer to render.
type FxKind uint8

const (
	FxPlaceRoot FxKind = iota
	FxPlaceWater
	FxPlaceMushroom
	FxPlaceCore
	FxHitCore
	FxDestroyRoot
	FxDestroyCore
	FxCutRot
	FxEatRot
)

type FxEvent struct {
	Kind     FxKind
	Col, Row int
}

func (g *Grid) emit(k FxKind, col, row int) {
	g.fxEvents = append(g.fxEvents, FxEvent{Kind: k, Col: col, Row: row})
}

// DrainFx returns the events emitted since the last call and clears the buffer.
func (g *Grid) DrainFx() []FxEvent {
	e := g.fxEvents
	g.fxEvents = nil
	return e
}
