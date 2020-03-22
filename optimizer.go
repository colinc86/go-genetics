package genetics

import (
	"math"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"github.com/cryptopirates/api"
	"github.com/cryptopirates/config"
)

// Optimizer types attempt to solve a problem given constraints.
type Optimizer struct {
	// A boolean that indicates whether or not the optimizer is running an optimization.
	Optimizing bool

	configuration   *config.OptimizerConfiguration
	evolver         *Evolver
	population      []*Chromosome
	chart           *api.Chart
	fitnessFunction func(chromosom *Chromosome, chart *api.Chart) float64
	generations     int
}

// MARK: Constructors

// NewOptimizerFromConfiguration creates and returns a new optimizer from the given configuration.
func NewOptimizerFromConfiguration(c *config.OptimizerConfiguration) *Optimizer {
	optimizer := &Optimizer{}

	evolverConfiguration := NewEvolverConfigurationFromOptimizerConfiguration(c)
	evolver := NewEvolver(evolverConfiguration, optimizer.optimizerFitnessFunction, optimizer.optimizerMutationFunction)

	optimizer.configuration = c
	optimizer.evolver = evolver
	optimizer.population = GeneratePopulation(uint(c.PopulationSize), uint(len(c.ChromosomeLimits)), optimizer.optimizerGenerationFunction)

	return optimizer
}

// MARK: Public functions

// Optimize optimizes solutions for the given chart.
func (o *Optimizer) Optimize(chart *api.Chart, fitnessFunction func(chromosom *Chromosome, chart *api.Chart) float64, finished func(c *Chromosome)) {
	o.Optimizing = true
	o.chart = chart
	o.fitnessFunction = fitnessFunction

	count := 0

	log.Debugf("Optimizer: running evolver for %d generations...", o.configuration.GenerationsPerCycle)
	o.evolver.Evolve(o.population, func(c *EvolverConfiguration, p Population) bool {
		count++
		o.generations++

		if count >= o.configuration.GenerationsPerCycle {
			o.Optimizing = false

			log.Debugf("Optimizer: finished optimizing parameters. (%d total generations.)\n", o.generations)
			finished(p[len(p)-1])
		}
		return count < o.configuration.GenerationsPerCycle
	})
}

// MARK: Private methods

func (o *Optimizer) optimizerFitnessFunction(chromosome *Chromosome) float64 {
	return o.fitnessFunction(chromosome, o.chart)
}

func (o *Optimizer) optimizerMutationFunction(chromosome *Chromosome, i int) float64 {
	value := chromosome.Genes[i]
	max := o.configuration.ChromosomeLimits[i].Max
	min := o.configuration.ChromosomeLimits[i].Min

	if o.configuration.ChromosomeLimits[i].Type == config.ChromosomeLimitTypeInt {
		return float64(rand.Intn(int(max)-int(min))) + float64(min)
	}

	step := math.Abs(max-min) / 10.0
	for j := 0; j < 3; j++ {
		if rand.Intn(2) == 1 {
			step /= 10.0
		}
	}

	if rand.Intn(2) == 1 {
		value += step
	} else {
		value -= step
	}

	if value > max {
		value = max
	} else if value < min {
		value = min
	}

	return value
}

func (o *Optimizer) optimizerGenerationFunction(i int, j int) float64 {
	max := o.configuration.ChromosomeLimits[j].Max
	min := o.configuration.ChromosomeLimits[j].Min

	if o.configuration.ChromosomeLimits[j].Type == config.ChromosomeLimitTypeInt {
		return float64(rand.Intn(int(max)-int(min))) + float64(min)
	}

	return rand.Float64()*(max-min) + min
}
