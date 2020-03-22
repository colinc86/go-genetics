package genetics

import (
	"strconv"
	"strings"

	"github.com/cryptopirates/config"
)

// EvolverConfiguration objects contains all of the necessary information needed
// to evolve a population of chromosomes using an evolver.
type EvolverConfiguration struct {
	SelectionMethod *SelectionMethod
	CrossoverMethod *CrossoverMethod
	Elitism         uint
	CrossoverRate   float64
	MutationRate    float64
}

// MARK: Constructors

// NewEvolverConfiguration creates and returns a new evolver configuration.
func NewEvolverConfiguration(selectionMethod *SelectionMethod, crossoverMethod *CrossoverMethod, elitism uint, crossoverRate float64, mutationRate float64) *EvolverConfiguration {
	return &EvolverConfiguration{
		SelectionMethod: selectionMethod,
		CrossoverMethod: crossoverMethod,
		Elitism:         elitism,
		CrossoverRate:   crossoverRate,
		MutationRate:    mutationRate,
	}
}

// NewEvolverConfigurationFromOptimizerConfiguration creates and returns a new evolver configuration.
func NewEvolverConfigurationFromOptimizerConfiguration(c *config.OptimizerConfiguration) *EvolverConfiguration {
	var selectionMethod *SelectionMethod
	switch c.SelectionMethod {
	case "rank":
		selectionMethod = NewSelectionMethod(SelectionMethodTypeRank)
	case "roulette":
		selectionMethod = NewSelectionMethod(SelectionMethodTypeRoulette)
	case "tournament":
		selectionMethod = NewSelectionMethod(SelectionMethodTypeTournament)
	}

	var crossoverMethod *CrossoverMethod
	if c.CrossoverMethod == "uniform" {
		crossoverMethod = NewCrossoverMethod(CrossoverMethodTypeUniform, 0)
	} else if strings.HasSuffix(c.CrossoverMethod, "-point") {
		components := strings.Split(c.CrossoverMethod, "-")
		if len(components) > 1 {
			if val, err := strconv.ParseInt(components[0], 10, 64); err == nil {
				crossoverMethod = NewCrossoverMethod(CrossoverMethodTypePoint, int(val))
			}
		}
	}

	config := NewEvolverConfiguration(selectionMethod, crossoverMethod, uint(c.Elitism), c.CrossoverRate, c.MutationRate)
	return config
}
