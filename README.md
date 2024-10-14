# blursh

This is a Go port for the Blurhash algorithm. To find out more about BlurHash, see https://github.com/woltapp/blurhash.

## Installation
Use go get

```
go get github.com/zar4za/blursh
```

## Usage

Create blurhash from image file
```golang
package main

import (
	"image"
	_ "image/jpeg"
	"os"

	"github.com/zar4za/blursh"
)

func main() {
	// please check for errors and close the file
	file, _ := os.Open("image.jpg")
	img, _, _ := image.Decode(file)

	hash, err := blursh.Encode(img, 4, 3)
}
```

`yComp` and `xComp` parameters adjust the amount of
vertical and horizontal AC components in hashed image. Both parameters must
be `>= 1` and `<= 8`. Basically it means how detailed the hash will be.
