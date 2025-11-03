package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/Nadim147c/material/v2"
)

func main() {
	// use ffmpeg command to covnert a media into raw rgb bytes
	cmd := exec.Command(
		"ffmpeg",
		"-i", "quantizer/gophar.jpg",
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-",
	)

	// returns bytes as r1,g1,b1,r2,g2,b2,r3,g3,b3
	rawRgbBytes, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	colors, err := material.Generate(
		material.FromBytes(rawRgbBytes),
		material.WithDark(true),
		material.WithVariant(material.VariantTonalSpot),
	)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(colors)
}
