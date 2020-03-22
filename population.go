package genetics

import (
	"math"
	"math/rand"
)

// Population types are an array of chromosomes.
type Population []*Chromosome

// MARK: Global methods

// GeneratePopulation generates a new population of chromosomes.
func GeneratePopulation(populationSize uint, chromosomeLength uint, generatingFunction func(i, j int) float64) Population {
	var population Population
	for i := 0; i < int(populationSize); i++ {
		chromosome := &Chromosome{}
		for j := 0; j < int(chromosomeLength); j++ {
			chromosome.Genes = append(chromosome.Genes, generatingFunction(i, j))
		}
		population = append(population, chromosome)
	}
	return population
}

// MARK: Public methods

// SumWeights returns the sum of the weights of the chromosomes in the population.
func (p Population) SumWeights() float64 {
	sum := 0.0
	for _, c := range p {
		sum += c.weight
	}
	return sum
}

// SumFitnesses returns the sum of the fitnesses of the chromosomes in the population.
func (p Population) SumFitnesses() float64 {
	sum := 0.0
	for _, c := range p {
		sum += c.Fitness
	}
	return sum
}

// CountNegativeWeights returns the number of chromosomes with negative weights
// in the population.
func (p Population) CountNegativeWeights() int {
	count := 0
	for _, c := range p {
		if c.weight < 0.0 {
			count++
		}
	}
	return count
}

// MinWeight returns the minimum weight of all chromosome in the population.
func (p Population) MinWeight() float64 {
	min := math.MaxFloat64
	for _, c := range p {
		if c.weight < min {
			min = c.weight
		}
	}
	return min
}

// ShiftWeights shifts the weights of the chromosomes of a population by a given
// value.
func (p Population) ShiftWeights(value float64) {
	for _, c := range p {
		c.weight += value
	}
}

// ShuffleChromosomes shuffles the chromosomes of the population.
func (p Population) ShuffleChromosomes() {
	rand.Shuffle(len(p), func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})
}

// ChromosomeWithMaxWeight returns the chromosome with the max weight in the population.
func (p Population) ChromosomeWithMaxWeight() *Chromosome {
	maxValue := -math.MaxFloat64
	maxIndex := 0
	for i, c := range p {
		if c.weight > maxValue {
			maxValue = c.weight
			maxIndex = i
		}
	}
	return p[maxIndex]
}
