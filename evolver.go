// Package genetics defines the objects necessary to perform a genetic evolution
// to solve complex problems.
package genetics

import (
	"math/rand"
	"sort"

	log "github.com/sirupsen/logrus"
)

// FitnessFunction defines a fitness function.
type FitnessFunction func(chromosome *Chromosome) float64

// MutationFunction defines a mutation function.
type MutationFunction func(chromosome *Chromosome, i int) float64

// Evolver types evolve a population given a configuration, fitness function,
// and mutation function.
type Evolver struct {
	Configuration    *EvolverConfiguration
	FitnessFunction  FitnessFunction
	MutationFunction MutationFunction
}

// MARK: Constructors

// NewEvolver creates and returns a new evolver.
func NewEvolver(configuration *EvolverConfiguration, fitnessFunction FitnessFunction, mutationFunction MutationFunction) *Evolver {
	return &Evolver{
		Configuration:    configuration,
		FitnessFunction:  fitnessFunction,
		MutationFunction: mutationFunction,
	}
}

// MARK: Public methods

// Evolve evolves a population.
func (e Evolver) Evolve(population Population, shouldContinue func(configuration *EvolverConfiguration, pop Population) bool) {
	if len(population) == 0 {
		log.Errorln("There are no chromosomes in the population.")
	}

	if e.Configuration.CrossoverMethod.Count >= len(population) {
		log.Errorln("The crossover count must be less than the number of chromosomes in the population.")
	}

	if int(e.Configuration.Elitism) > len(population) {
		log.Errorln("The elitism count must be less than or equal to the number of chromosomes in the population.")
	}

	e.calculateFitnesses(population)
	sort.Slice(population[:], func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	for shouldContinue(e.Configuration, population) {
		population = e.breedSingleGeneration(population)
		e.calculateFitnesses(population)

		sort.Slice(population[:], func(i, j int) bool {
			return population[i].Fitness < population[j].Fitness
		})
	}
}

// MARK: Private methods

// shouldCrossover returns whether or not the evolver should perform crossover.
func (e Evolver) shouldCrossover() bool {
	return rand.Float64() <= e.Configuration.CrossoverRate
}

// shouldMutate returns whether or not the evolver should perform mutation.
func (e Evolver) shouldMutate() bool {
	return rand.Float64() <= e.Configuration.MutationRate
}

// calculateFitness calculates the fitness of each chromosome in a population.
func (e Evolver) calculateFitnesses(population Population) {
	for i := 0; i < len(population); i++ {
		fitness := e.FitnessFunction(population[i])
		if fitness < 0.0 {
			// log.Warnf("Negative fitness value %f may cause strange results.", fitness)
		}

		population[i].Fitness = fitness
		population[i].weight = fitness
	}
}

// breedSingleGeneration breeds a single generation of chromosomes from a population.
func (e Evolver) breedSingleGeneration(population Population) Population {
	var newPopulation Population
	elite := e.applyElitism(population)

	newPopulation = append(newPopulation, elite...)

	for i := len(elite); i < len(population); i++ {
		child := e.breedChild(population)
		// log.Debugf("Got child %s\n", child)
		newPopulation = append(newPopulation, child)
	}

	return newPopulation
}

// applyElitisim applies elitism to a population and places the chromosomes that
// survived in to the destination population.
func (e Evolver) applyElitism(population Population) []*Chromosome {
	var chromosomes []*Chromosome
	for i := 0; i < int(e.Configuration.Elitism); i++ {
		chromosomes = append(chromosomes, population[len(population)-i-1])
	}
	return chromosomes
}

// breedChild breeds a child chromosome from the population.
func (e Evolver) breedChild(population Population) *Chromosome {
	child := &Chromosome{}
	child.Genes = make([]float64, len(population[0].Genes))

	if e.shouldCrossover() {
		chromosome := e.Configuration.CrossoverMethod.Function(
			e.Configuration.SelectionMethod.Function(population),
			e.Configuration.SelectionMethod.Function(population),
			e.Configuration.CrossoverMethod.Count,
		)
		copy(child.Genes, chromosome.Genes)
		child.Fitness = chromosome.Fitness
		child.weight = chromosome.weight
	} else {
		chromosome := e.Configuration.SelectionMethod.Function(population)
		copy(child.Genes, chromosome.Genes)
		child.Fitness = chromosome.Fitness
		child.weight = chromosome.weight
	}

	for i := 0; i < len(child.Genes); i++ {
		if e.shouldMutate() {
			child.Genes[i] = e.MutationFunction(child, i)
		}
	}
	// log.Debugf("Returning child %s\n", child)
	return child
}
