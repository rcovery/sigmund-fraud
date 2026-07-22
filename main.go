package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

		vectors := vectorize(&parsedBody)
		_, writeErr := fmt.Fprintf(w, "%v", vectors)
		if writeErr != nil {
			log.Fatalln(writeErr)
		}

		c := GetCentroid(vectors)
		restoredVectors := LoadFromIVF(c)

		for _, restoredVector := range *restoredVectors {
			distance := EuclideanDistance(vectors, restoredVector.Vector)
			log.Printf("Distance: %v", distance)
		}
	})

	srvErr := srv.ListenAndServe()
	if srvErr != nil {
		log.Fatalln(srvErr)
	}
}
