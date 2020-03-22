package genetics

import (
	"math/rand"
	"sort"
)

// SelectionMethodType represents a type of selection method.
type SelectionMethodType uint

// Types of selection methods.
const (
	SelectionMethodTypeRank       SelectionMethodType = 0
	SelectionMethodTypeRoulette   SelectionMethodType = 1
	SelectionMethodTypeTournament SelectionMethodType = 2
	SelectionMethodTypeCustom     SelectionMethodType = 3
)

// SelectionMethodFunction takes a population of chromosomes and chooses one for
// breeding.
type SelectionMethodFunction func(population Population) *Chromosome

// SelectionMethod wraps a method type and function together.
type SelectionMethod struct {
	Type     SelectionMethodType
	Function SelectionMethodFunction
}

// MARK: Constructors

// NewSelectionMethod creates a new selection method from the given selection
// method type. To use a custom function, use the `NewCustomSelectionMethod`
// constructor.
func NewSelectionMethod(t SelectionMethodType) *SelectionMethod {
	return &SelectionMethod{
		Type:     t,
		Function: selectionFunctionForType(t),
	}
}

// NewCustomSelectionMethod creates a new custom selection method from the
// provided selection method function.
func NewCustomSelectionMethod(f SelectionMethodFunction) *SelectionMethod {
	return &SelectionMethod{
		Type:     SelectionMethodTypeCustom,
		Function: f,
	}
}

// MARK: Public functions

// RankFunction implements the rank selection function.
var RankFunction SelectionMethodFunction = func(population Population) *Chromosome {
	for i := 0; i < len(population); i++ {
		population[i].weight = float64(i) + 1.0
	}

	total := population.SumWeights()
	rand := float64(rand.Intn(int(total)))
	sum := 0.0

	for _, c := range population {
		sum += c.weight
		if rand < sum {
			return c
		}
	}

	return &Chromosome{}
}

// RouletteFunction implements the roulette selection function.
var RouletteFunction SelectionMethodFunction = func(population Population) *Chromosome {
	sort.Slice(population[:], func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})

	total := population.SumWeights()
	for i := 0; i < len(population); i++ {
		population[i].weight = population[i].weight / total
	}

	count := population.CountNegativeWeights()
	if count > 0 {
		// log.Warnf("Population contains %d chromosomes with fitness < 0.0 which may result in bad roulette selection.", count)
		min := population.MinWeight()
		population.ShiftWeights(min)
	}

	rand := rand.Float64()
	sum := 0.0

	for _, c := range population {
		sum += c.weight
		if rand < sum {
			return c
		}
	}

	return population[0]
}

// TournamentFunction implements the tournament selection function.
var TournamentFunction SelectionMethodFunction = func(population Population) *Chromosome {
	population.ShuffleChromosomes()
	rand := rand.Intn(len(population)-1) + 1
	tournamentGroup := population[0:rand]
	return tournamentGroup.ChromosomeWithMaxWeight()
}

// MARK: Private functions

// selectionFunctionForType returns the selection function for the given type.
func selectionFunctionForType(t SelectionMethodType) SelectionMethodFunction {
	switch t {
	case SelectionMethodTypeRank:
		return RankFunction
	case SelectionMethodTypeRoulette:
		return RouletteFunction
	case SelectionMethodTypeTournament:
		return TournamentFunction
	default:
		return nil
	}
}
