package main

import (
	"log"
	"math"
)

func main() {
	// Try to use govector
	exampleDataset := make([]Vector, 0)
	exampleDataset = append(exampleDataset,
		Vector{
			[]float64{0.9708, 1.0000, 1.00},
			"fraud",
		})

	dimensionAmount := clamp(12000.0 / 10000.0)
	dimensionHour := clamp(22.0 / 23.0)
	dimensionAvg := clamp(4800.0 / 5000.0)

	myvector := []float64{dimensionAmount, dimensionHour, dimensionAvg}

	log.Print(euclideanDistance(myvector, exampleDataset[0].vector))
}

func euclideanDistance(myvec []float64, dataset []float64) float64 {
	var totalDiff float64
	for idx, dim := range dataset {
		diff := dim - myvec[idx]
		totalDiff += math.Pow(diff, 2.0)
	}
	return math.Sqrt(totalDiff)
}

func clamp(data float64) float64 {
	if data < 0 {
		return 0
	}
	if data > 1 {
		return 1
	}

	return data
}

type Vector struct {
	vector []float64
	label  TxType
}

type TxType string

const (
	TxLegit TxType = "legit"
	TxFraud TxType = "fraud"
)
