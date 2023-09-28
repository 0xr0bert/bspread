package beliefs

import (
	"github.com/0xr0bert/bspread/behaviours"
	"github.com/google/uuid"
)

type Belief struct {
	Name                    string
	Uuid                    uuid.UUID
	Relationships           map[*Belief]float64
	Perceptions             map[*behaviours.Behaviour]float64
	PerformanceRelationship map[*behaviours.Behaviour]float64
}

func New(name string) (b *Belief) {
	b = new(Belief)
	b.Name = name
	b.Uuid = uuid.New()
	b.Relationships = make(map[*Belief]float64)
	b.Perceptions = make(map[*behaviours.Behaviour]float64)
	b.PerformanceRelationship = make(map[*behaviours.Behaviour]float64)

	return
}
