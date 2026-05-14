package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	MaxAmount            float64 = 10000
	MaxInstallments      float64 = 12
	AmountVsAvgRatio     float64 = 10
	MaxMinutes           float64 = 1440
	MaxKm                float64 = 1000
	MaxTxCount24h        float64 = 20
	MaxMerchantAvgAmount float64 = 10000
)

var MccRisk MccRiskPayload = MccRiskPayload{
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

func main() {
	srv := http.Server{
		Addr: ":9999",
	}

	http.HandleFunc("/fraud-score", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			return
		}

		rawBody, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			w.WriteHeader(400)
			_, writeErr := fmt.Fprintf(w, "Body inválido: %v", bodyErr)
			if writeErr != nil {
				log.Fatalln(writeErr)
			}
			return
		}

		var parsedBody FraudScoreBody
		jsonErr := json.Unmarshal(rawBody, &parsedBody)
		if jsonErr != nil {
			w.WriteHeader(400)
			_, writeErr := fmt.Fprintf(w, "Json inválido: %v", jsonErr)
			if writeErr != nil {
				log.Fatalln(writeErr)
			}
			return
		}

		dimensions := vectorize(&parsedBody)
		_, writeErr := fmt.Fprintf(w, "%v", dimensions)
		if writeErr != nil {
			log.Fatalln(writeErr)
		}
	})

	srvErr := srv.ListenAndServe()
	if srvErr != nil {
		log.Fatalln(srvErr)
	}
}

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
