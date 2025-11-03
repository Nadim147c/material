package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Nadim147c/material/v2"
)

func main() {
	colors, err := material.Generate(
		material.FromHex("#0044ff"), // genearte colors from a blue color
		material.WithContrast(0.3),
		material.WithVariant(material.VariantContent),
		material.WithVersion(material.Version2025),
	)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(colors)
}
