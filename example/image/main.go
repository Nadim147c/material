package main

import (
	"encoding/json"
	"image/jpeg"
	"log"
	"os"

	"github.com/Nadim147c/material/v2"
)

func main() {
	f, err := os.Open("quantizer/gophar.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	colors, err := material.Generate(
		material.FromImage(img),
		material.WithDark(true),
		material.WithVariant(material.VariantTonalSpot),
	)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(colors)
}
