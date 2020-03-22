package genetics

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
