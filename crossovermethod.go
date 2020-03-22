package genetics

import (
	"math/rand"
	"sort"
)

// CrossoverMethodType represents a type of crossover method.
type CrossoverMethodType uint

// Types of crossover methods.
const (
	CrossoverMethodTypePoint   CrossoverMethodType = 0
	CrossoverMethodTypeUniform CrossoverMethodType = 1
	CrossoverMethodTypeCustom  CrossoverMethodType = 2
)

// CrossoverMethodFunction takes a pair of chromosomes and performs crossover
// between them.
type CrossoverMethodFunction func(cA *Chromosome, cB *Chromosome, count int) *Chromosome

// CrossoverMethod wraps a method type and function together.
type CrossoverMethod struct {
	Type     CrossoverMethodType
	Function CrossoverMethodFunction
	Count    int
}

// MARK: Constructors

// NewCrossoverMethod creates a new crossover method from the given crossover
// method type. To use a custom function, use the `NewCustomCrossoverMethod`
// constructor.
func NewCrossoverMethod(t CrossoverMethodType, count int) *CrossoverMethod {
	return &CrossoverMethod{
		Type:     t,
		Function: crossoverFunctionForType(t),
		Count:    count,
	}
}

// NewCustomCrossoverMethod creates a new custom crossover method from the
// provided crossover method function.
func NewCustomCrossoverMethod(f CrossoverMethodFunction, count int) *CrossoverMethod {
	return &CrossoverMethod{
		Type:     CrossoverMethodTypeCustom,
		Function: f,
		Count:    count,
	}
}

// MARK: Public functions

// PointFunction implements the point crossover function.
var PointFunction CrossoverMethodFunction = func(cA *Chromosome, cB *Chromosome, count int) *Chromosome {
	var indexes []int
	for i := 0; i < len(cA.Genes); i++ {
		indexes = append(indexes, i+1)
	}

	if len(indexes) > 1 {
		for i := 0; i < len(indexes)-1; i++ {
			j := rand.Intn(len(indexes)-i) + i
			if i == j {
				continue
			}
			indexes[i], indexes[j] = indexes[j], indexes[i]
		}
	}

	crossoverPoints := indexes[0:count]
	sort.Ints(crossoverPoints)
	crossoverPoints = append([]int{0}, crossoverPoints...)
	crossoverPoints = append(crossoverPoints, len(cA.Genes))

	child := &Chromosome{}
	for _, g := range cA.Genes {
		child.Genes = append(child.Genes, g)
	}

	for i := 0; i < len(crossoverPoints)-1; i++ {
		for j := crossoverPoints[i]; j < crossoverPoints[i+1]; j++ {
			if i%2 == 0 {
				child.Genes[j] = cA.Genes[j]
			} else {
				child.Genes[j] = cB.Genes[j]
			}
		}
	}

	return child
}

// UniformFunction implements the uniform crossover function.
var UniformFunction CrossoverMethodFunction = func(cA *Chromosome, cB *Chromosome, count int) *Chromosome {
	child := &Chromosome{}
	for _, g := range cA.Genes {
		child.Genes = append(child.Genes, g)
	}

	for i := 0; i < len(cA.Genes); i++ {
		if rand.Intn(2) == 1 {
			child.Genes[i] = cA.Genes[i]
		} else {
			child.Genes[i] = cB.Genes[i]
		}
	}

	return child
}

// MARK: Private functions

// crossoverFunctionForType returns the crossover function for the given type.
func crossoverFunctionForType(t CrossoverMethodType) CrossoverMethodFunction {
	switch t {
	case CrossoverMethodTypePoint:
		return PointFunction
	case CrossoverMethodTypeUniform:
		return UniformFunction
	default:
		return nil
	}
}
