# Material

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Nadim147c/material?style=for-the-badge&logo=go&labelColor=11140F&color=BBE9AA)](https://pkg.go.dev/github.com/Nadim147c/material)
[![GitHub Repo stars](https://img.shields.io/github/stars/Nadim147c/material?style=for-the-badge&logo=github&labelColor=11140F&color=BBE9AA)](https://github.com/Nadim147c/material)
[![GitHub License](https://img.shields.io/github/license/Nadim147c/material?style=for-the-badge&logo=gplv3&labelColor=11140F&color=BBE9AA)](./LICENSE)
[![GitHub Tag](https://img.shields.io/github/v/tag/Nadim147c/material?include_prereleases&sort=semver&style=for-the-badge&logo=git&labelColor=11140F&color=BBE9AA)](https://github.com/Nadim147c/material/tags)

> [!IMPORTANT]
> 🔥 Found this useful? A quick star goes a long way.

A pure go implementation of [Material Color
Utilities](https://github.com/material-foundation/material-color-utilities) without
any external dependencies.

## Example

[Example](./example).

```go
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

	json.NewEncoder(os.Stdout).Encode(colors) // prints the colors as json
}
```

## License

This project is licensed under the [Apache License, Version 2.0](./LICENSE).

It incorporates code derived from [Material Color
Utilities](https://github.com/material-foundation/material-color-utilities) by
**Google LLC**, used under the terms of the Apache License, Version 2.0.

The included [gophar image](./quantizer/gophar.jpg), used as test data, is licensed
under the **Creative Commons Attribution 4.0 International (CC BY 4.0)** license,
attributed to **Renee French**.

For additional details and third-party attributions, see the [NOTICE](./NOTICE) file.
