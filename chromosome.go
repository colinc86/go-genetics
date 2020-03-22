package genetics

import "fmt"

// Chromosome object contain an array of genes and a fitness value.
type Chromosome struct {
	// The chromosome's genes.
	Genes []float64

	// The fitness of the chromosome. If the chromosome is part of a `Population`
	// object, then this value is updated each time the population evolves. To
	// prevent excessive calls to the `Evolver`'s `FitnessFunction`, this value is
	// only updated _once_ immediately following evolution for each generation.
	// Changing the value of this property before the next evolution *will* affect
	// selection from the population.
	Fitness float64

	// The weight of the chromosome. Internal use only.
	weight float64
}

// MARK: String methods

func (c Chromosome) String() string {
	return fmt.Sprintf("[Genes: %v, Fitness: %0.10f, weight: %0.10f]", c.Genes, c.Fitness, c.weight)
}
