package main

var (
	MaxAmount            float64 = 10000
	MaxInstallments      float64 = 12
	AmountVsAvgRatio     float64 = 10
	MaxMinutes           float64 = 1440
	MaxKm                float64 = 1000
	MaxTxCount24h        float64 = 20
	MaxMerchantAvgAmount float64 = 10000
)

var MccRisk = MccRiskPayload{
	"5411": 0.15,
	"5812": 0.30,
	"5912": 0.20,
	"5944": 0.45,
	"7801": 0.80,
	"7802": 0.75,
	"7995": 0.85,
	"4511": 0.35,
	"5311": 0.25,
	"5999": 0.50,
}
