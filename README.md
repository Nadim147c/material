# Material

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Nadim147c/material?style=for-the-badge&logo=go&labelColor=11140F&color=BBE9AA)](https://pkg.go.dev/github.com/Nadim147c/material)
[![GitHub Repo stars](https://img.shields.io/github/stars/Nadim147c/material?style=for-the-badge&logo=github&labelColor=11140F&color=BBE9AA)](https://github.com/Nadim147c/material)
[![GitHub License](https://img.shields.io/github/license/Nadim147c/material?style=for-the-badge&logo=gplv3&labelColor=11140F&color=BBE9AA)](./LICENSE)
[![GitHub Tag](https://img.shields.io/github/v/tag/Nadim147c/material?include_prereleases&sort=semver&style=for-the-badge&logo=git&labelColor=11140F&color=BBE9AA)](https://github.com/Nadim147c/material/tags)

> [!IMPORTANT]
> 🔥 Found this useful? A quick star goes a long way.

A pure go implementation of [Material Color
Utilities](https://github.com/material-foundation/material-color-utilities)

## Example

```go
package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/Nadim147c/material"
	"github.com/Nadim147c/material/dynamic"
)

func main() {
	file, err := os.Open("gophar.jpg")
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	colors, err := material.GenerateFromImage(
		img,
		dynamic.VariantExpressive,
		true,
		0,
		dynamic.PlatformPhone,
		dynamic.Version2021,
	)
	if err != nil {
		log.Fatalf("failed to generate colors: %v", err)
	}

	for key, value := range colors {
		fmt.Println(key, value)
	}
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
