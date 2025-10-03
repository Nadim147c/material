# Material

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Nadim147c/material?style=for-the-badge&logo=go&labelColor=11140F&color=BBE9AA)](https://pkg.go.dev/github.com/Nadim147c/material)
[![GitHub Repo stars](https://img.shields.io/github/stars/Nadim147c/material?style=for-the-badge&logo=github&labelColor=11140F&color=BBE9AA)](https://github.com/Nadim147c/material)
[![GitHub License](https://img.shields.io/github/license/Nadim147c/material?style=for-the-badge&logo=gplv3&labelColor=11140F&color=BBE9AA)](./LICENSE)
[![GitHub Tag](https://img.shields.io/github/v/tag/Nadim147c/material?include_prereleases&sort=semver&style=for-the-badge&logo=git&labelColor=11140F&color=BBE9AA)](https://github.com/Nadim147c/material/tags)

> [!IMPORTANT]
> 🔥 Found this useful? A quick star goes a long way.

> [!CAUTION]
> This is tool is in beta and expect some unforeseen bugs.

Pure go implementation of [Material Color Utilities](https://github.com/material-foundation/material-color-utilities)

## Example

```go
package main

import (
	"github.com/Nadim147c/material"
	"github.com/Nadim147c/material/dynamic"
)

func main() {
	file, err := gophar.Open("gophar.jpg")
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("failed to decode image: %v", err)
	}

	colors, err := material.GenerateFromImage(img, dynamic.Expressive, true, 0, dynamic.Phone, dynamic.V2021)
	if err != nil {
		t.Fatalf("failed to generate colors: %v", err)
	}

	for key, value := range colors {
		t.Log(key, value)
	}
}
```

## License

This project is licensed under the [Apache License, Version 2.0](./LICENSE). It
includes code derived from
[Material Color Utilities](https://github.com/material-foundation/material-color-utilities)
by Google LLC, originally licensed under the
[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

See the [NOTICE](./NOTICE) file for details and third-party attributions.
