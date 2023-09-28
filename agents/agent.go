package agents

import (
	"github.com/0xr0bert/bspread/behaviours"
	"github.com/google/uuid"
	"math/rand"
	"sort"
)
import "github.com/0xr0bert/bspread/beliefs"
import "github.com/0xr0bert/bspread/simulation"

type Agent struct {
	Uuid            uuid.UUID
	Activation      map[simulation.SimTime]map[*beliefs.Belief]float64
	Friends         map[*Agent]float64
	Actions         map[simulation.SimTime]*behaviours.Behaviour
	Deltas          map[*beliefs.Belief]float64
	KnownBehaviours []*behaviours.Behaviour
}

func (a *Agent) GetActivation(time simulation.SimTime, belief *beliefs.Belief) (float64, bool) {
	if a.Activation == nil {
		return 0, false
	}

	innerMap := a.Activation[time]

	if innerMap == nil {
		return 0, false
	} else {
		return innerMap[belief], true
	}
}

func (a *Agent) SetActivation(time simulation.SimTime, belief *beliefs.Belief, activation float64) {
	if a.Activation == nil {
		a.Activation = make(map[simulation.SimTime]map[*beliefs.Belief]float64)
	}

	if a.Activation[time] == nil {
		a.Activation[time] = make(map[*beliefs.Belief]float64)
	}

	a.Activation[time][belief] = activation
}

func (a *Agent) WeightedRelationship(b1 *beliefs.Belief, b2 *beliefs.Belief, time simulation.SimTime) float64 {
	act := a.Activation[time][b1]
	r := b1.Relationships[b2]

	return act * r
}

func (a *Agent) Contextualise(b *beliefs.Belief, time simulation.SimTime) float64 {
	activations := a.Activation[time]

	if activations == nil {
		return 0
	}

	n := 0
	sum := 0.0

	for b2 := range activations {
		sum += a.WeightedRelationship(b, b2, time)
		n++
	}

	return sum / float64(n)
}

func (a *Agent) Pressure(b *beliefs.Belief, time simulation.SimTime) float64 {
	n := 0
	sum := 0.0

	for friend, weight := range a.Friends {
		sum += weight * b.Perceptions[friend.Actions[time]]
		n++
	}

	return sum / float64(n)
}

func (a *Agent) ActivationChange(b *beliefs.Belief, time simulation.SimTime) float64 {
	pressure := a.Pressure(b, time)
	context := a.Contextualise(b, time)

	var weight float64

	if pressure > 0 {
		weight = (1.0 + context) / 2.0
	} else {
		weight = (1.0 - context) / 2.0
	}

	return weight * pressure
}

func (a *Agent) updateActivation(b *beliefs.Belief, time simulation.SimTime) {
	a.SetActivation(time, b,
		max(-1, min(1, a.Deltas[b]*a.Activation[time-1][b]+a.ActivationChange(b, time-1))),
	)
}

func (a *Agent) updateAllActivations(time simulation.SimTime) {
	for b := range a.Activation[time-1] {
		a.updateActivation(b, time)
	}
}

type probPair struct {
	behaviour *behaviours.Behaviour
	value     float64
}

func (a *Agent) calculateSortedUnnormalizedProbs(time simulation.SimTime) (probs []probPair) {
	probs = make([]probPair, len(a.KnownBehaviours))

	for i, behaviour := range a.KnownBehaviours {
		probs[i].behaviour = behaviour
		value := 0.0

		for belief, act := range a.Activation[time] {
			prs := belief.PerformanceRelationship[behaviour]
			value += prs * act
		}

		probs[i].value = value
	}

	sort.Slice(probs, func(i, j int) bool {
		return probs[i].value < probs[j].value
	})

	return
}

func (a *Agent) chooseActionIfNoneHavePositiveProbs(probs []probPair, time simulation.SimTime) {
	a.Actions[time] = probs[len(probs)-1].behaviour
}

func (a *Agent) chooseActionIfOnlyOneIsPositive(filteredProbs []probPair, time simulation.SimTime) {
	a.Actions[time] = filteredProbs[0].behaviour
}

func (a *Agent) chooseActionIfMoreThanOneIsPositive(filteredProbs []probPair, time simulation.SimTime) {
	normalizingFactor := 0.0

	for _, p := range filteredProbs {
		normalizingFactor += p.value
	}

	normalizedProbs := make([]probPair, len(filteredProbs))

	for i, p := range filteredProbs {
		normalizedProbs[i].behaviour = p.behaviour
		normalizedProbs[i].value = p.value / normalizingFactor
	}

	chosenBehaviour := normalizedProbs[len(normalizedProbs)-1].behaviour

	rv := rand.Float64()

	for _, p := range normalizedProbs {
		rv -= p.value
		if rv <= 0 {
			chosenBehaviour = p.behaviour
			break
		}
	}

	a.Actions[time] = chosenBehaviour
}

func filterPositiveProbs(probs []probPair) []probPair {
	var filteredProbs []probPair
	for _, p := range probs {
		if p.value >= 0 {
			filteredProbs = append(filteredProbs, p)
		}
	}

	return filteredProbs
}

func (a *Agent) chooseAction(time simulation.SimTime) {
	probs := a.calculateSortedUnnormalizedProbs(time)

	if probs[len(probs)-1].value < 0 {
		a.chooseActionIfNoneHavePositiveProbs(probs, time)
	} else {
		filteredProbs := filterPositiveProbs(probs)

		if len(filteredProbs) == 1 {
			a.chooseActionIfOnlyOneIsPositive(filteredProbs, time)
		} else {
			a.chooseActionIfMoreThanOneIsPositive(filteredProbs, time)
		}
	}
}
