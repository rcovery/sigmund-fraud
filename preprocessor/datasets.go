// Package preprocessor serve para preprocessar os datasets e salvar em gob
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func datasets() {
	log.Print("test")
	f, _ := os.ReadFile("processor/references.json")
	jsonLoaded := json.NewDecoder(bytes.NewReader(f))

	var jsonDecoded []Vector
	jsonLoaded.Decode(&jsonDecoded)

	for _, jsonEntry := range jsonDecoded {
		var centroid float64
		for _, vectorEntry := range jsonEntry.Vector {
			centroid += vectorEntry
		}
		centroid /= float64(len(jsonEntry.Vector))

		os.WriteFile(fmt.Sprintf("processor/%.2f", centroid), []byte(fmt.Sprint(jsonEntry.Vector)), 0o644)
		break
	}
}

type Vector struct {
	Vector []float64
	Label  string
}
