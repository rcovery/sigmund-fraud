package main

import "time"

type Transaction struct {
	RequestedAt  time.Time `json:"requested_at"`
	Amount       float64   `json:"amount"`
	Installments float64   `json:"installments"`
	_            [4]byte
}

type Customer struct {
	KnownMerchants []string `json:"known_merchants"`
	AvgAmount      float64  `json:"avg_amount"`
	TxCount24h     int32    `json:"tx_count_24h"`
	_              [4]byte
}

type Merchant struct {
	ID        string  `json:"id"`
	MCC       string  `json:"mcc"`
	AvgAmount float64 `json:"avg_amount"`
}

type Terminal struct {
	KmFromHome  float64 `json:"km_from_home"`
	IsOnline    bool    `json:"is_online"`
	CardPresent bool    `json:"card_present"`
	_           [6]byte
}

type LastTransaction struct {
	Timestamp     time.Time `json:"timestamp"`
	KmFromCurrent float64   `json:"km_from_current"`
}

type FraudScoreBody struct {
	Transaction     Transaction      `json:"transaction"`
	Customer        Customer         `json:"customer"`
	LastTransaction *LastTransaction `json:"last_transaction"`
	Merchant        Merchant         `json:"merchant"`
	ID              string           `json:"id"`
	Terminal        Terminal         `json:"terminal"`
}

type FraudScoreResponse struct {
	FraudScore float64 `json:"fraud_score"`
	Approved   bool    `json:"approved"`
	_          [7]byte
}

type MccRiskPayload map[string]float64
