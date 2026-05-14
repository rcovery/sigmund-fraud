package main

import (
	"math"
	"slices"
	"time"
)

func vectorize(body *FraudScoreBody) []float64 {
	/*
		0	amount	limitar(transaction.amount / max_amount)
		1	installments	limitar(transaction.installments / max_installments)
		2	amount_vs_avg	limitar((transaction.amount / customer.avg_amount) / amount_vs_avg_ratio)
		3	hour_of_day	hora(transaction.requested_at) / 23 (0-23, UTC)
		4	day_of_week	dia_da_semana(transaction.requested_at) / 6 (seg=0, dom=6)
		5	minutes_since_last_tx	limitar(minutos / max_minutes) ou -1 se last_transaction: null
		6	km_from_last_tx	limitar(last_transaction.km_from_current / max_km) ou -1 se last_transaction: null
		7	km_from_home	limitar(terminal.km_from_home / max_km)
		8	tx_count_24h	limitar(customer.tx_count_24h / max_tx_count_24h)
		9	is_online	1 se terminal.is_online, senão 0
		10	card_present	1 se terminal.card_present, senão 0
		11	unknown_merchant	1 se merchant.id não estiver em customer.known_merchants, senão 0 (invertido: 1 = desconhecido)
		12	mcc_risk	mcc_risk.json[merchant.mcc] (valor padrão 0.5)
		13	merchant_avg_amount	limitar(merchant.avg_amount / max_merchant_avg_amount)
	*/
	var dimensions []float64
	dimensions = append(dimensions, Clamp(body.Transaction.Amount/MaxAmount))
	dimensions = append(dimensions, Clamp(body.Transaction.Installments/MaxInstallments))
	dimensions = append(dimensions, Clamp((body.Transaction.Amount/body.Customer.AvgAmount)/AmountVsAvgRatio))
	dimensions = append(dimensions, float64(body.Transaction.RequestedAt.Hour())/23)
	dimensions = append(dimensions, GetVecWeekday(body.Transaction.RequestedAt))
	dimensions = append(dimensions, GetVecMinutesSinceLastTx(body.LastTransaction))
	dimensions = append(dimensions, GetVecKmFromLastTx(body.LastTransaction))
	dimensions = append(dimensions, body.Terminal.KmFromHome/MaxKm)
	dimensions = append(dimensions, float64(body.Customer.TxCount24h)/MaxTxCount24h)
	dimensions = append(dimensions, GetVecBool(body.Terminal.IsOnline))
	dimensions = append(dimensions, GetVecBool(body.Terminal.CardPresent))
	dimensions = append(dimensions, GetVecKnownMerchant(body))
	dimensions = append(dimensions, GetVecMccRisk(body.Merchant.ID))
	dimensions = append(dimensions, body.Merchant.AvgAmount/MaxMerchantAvgAmount)

	return dimensions
}

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
