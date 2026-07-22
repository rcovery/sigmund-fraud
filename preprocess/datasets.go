// Package preprocessor serve para preprocessar os datasets e salvar em gob
package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	f, _ := os.ReadFile("preprocess/references.json")
	jsonLoaded := json.NewDecoder(bytes.NewReader(f))

	var jsonDecoded []Vector
	jsonLoaded.Decode(&jsonDecoded)

	mappedCentroids := map[string][]Vector{}

	for _, jsonEntry := range jsonDecoded {
		var centroid float64
		for _, vectorEntry := range jsonEntry.Vector {
			centroid += vectorEntry
		}
		centroid /= float64(len(jsonEntry.Vector))
		if centroid < 0 {
			centroid = 0
		}

		stringifiedCentroid := fmt.Sprintf("%.3f", centroid)

		mappedCentroids[stringifiedCentroid] = append(mappedCentroids[stringifiedCentroid], jsonEntry)
	}

	for key, item := range mappedCentroids {
		gobFileName := fmt.Sprintf("preprocess/preprocessed-centroids/%v", key)
		log.Printf("Writing at %v", gobFileName)

		w, _ := os.Create(gobFileName)
		defer w.Close()

		enc := gob.NewEncoder(w)
		enc.Encode(item)
	}
}

type Vector struct {
	Vector []float64
	Label  string
}
