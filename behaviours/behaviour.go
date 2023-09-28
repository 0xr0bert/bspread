package behaviours

import "github.com/google/uuid"

type Behaviour struct {
	Name string
	Uuid uuid.UUID
}

func New(name string) (b *Behaviour) {
	b = new(Behaviour)
	b.Name = name
	b.Uuid = uuid.New()

	return
}
