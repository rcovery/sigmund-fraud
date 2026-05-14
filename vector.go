package main

import (
	"math"
	"slices"
	"time"
)

func GetVecMccRisk(merchantID string) float64 {
	risk, ok := MccRisk[merchantID]
	if !ok {
		return 0.5
	}

	return risk
}

func GetVecKnownMerchant(body *FraudScoreBody) float64 {
	if slices.Contains(body.Customer.KnownMerchants, body.Merchant.ID) {
		return 1
	}

	return 0
}

func GetVecBool(data bool) float64 {
	if data {
		return 1
	}

	return 0
}

func GetVecWeekday(t time.Time) float64 {
	weekday := (int(t.UTC().Weekday()) + 6) % 7
	return float64(weekday) / 6.0
}

func GetVecMinutesSinceLastTx(lastTx *LastTransaction) float64 {
	if lastTx == nil {
		return -1
	}

	return float64(lastTx.Timestamp.Minute()) / MaxMinutes
}

func GetVecKmFromLastTx(lastTx *LastTransaction) float64 {
	if lastTx == nil {
		return -1
	}

	return float64(lastTx.KmFromCurrent) / MaxKm
}

func EuclideanDistance(myvec []float64, dataset []float64) float64 {
	var totalDiff float64
	for idx, dim := range dataset {
		diff := dim - myvec[idx]
		totalDiff += math.Pow(diff, 2.0)
	}
	return math.Sqrt(totalDiff)
}

func Clamp(data float64) float64 {
	if data < 0 {
		return 0.0
	}
	if data > 1 {
		return 1.0
	}

	return data
}
