package blursh

import (
	"errors"
	"image"
	"image/color"
	"math"
	"strings"
	"sync"
)

type factor struct {
	r float64
	g float64
	b float64
}

type pixel struct {
	r, g, b uint8
}

func Encode(img image.Image, xComp int, yComp int) (string, error) {
	if xComp < 1 || xComp > 9 {
		return "", errors.New("xComp must be in range from 1 to 9")
	}
	if yComp < 1 || yComp > 9 {
		return "", errors.New("yComp must be in range from 1 to 9")
	}

	pixels, width, height := imageToPixels(img)
	factors := make([]factor, xComp*yComp)
	wg := sync.WaitGroup{}

	for y := range yComp {
		for x := range xComp {
			wg.Add(1)
			go multiplyBasisFunction(pixels, x, y, width, height, &factors[y*xComp+x], &wg)
		}
	}

	ac := factors[1:]
	builder := strings.Builder{}
	builder.Grow(4 + xComp*yComp*2)
	sizeFlag := (xComp - 1) + (yComp-1)*9
	encode83(&builder, sizeFlag, 1)

	wg.Wait()

	actualMaximumValue := 0.
	for _, factor := range ac {
		actualMaximumValue = math.Max(math.Max(math.Max(factor.r, factor.g), factor.b), actualMaximumValue)
	}

	quantisedMaximumValue := math.Max(0, math.Min(82, actualMaximumValue*166-0.5))
	maximumValue := (quantisedMaximumValue + 1) / 166
	encode83(&builder, int(quantisedMaximumValue), 1)
	encode83(&builder, encodeDC(factors[0]), 4)

	for _, factor := range ac {
		encode83(&builder, encodeAC(factor, maximumValue), 2)
	}

	return builder.String(), nil
}

func multiplyBasisFunction(pixels []pixel, xComp, yComp, width, height int, fct *factor, wg *sync.WaitGroup) {
	defer wg.Done()
	result := factor{}
	normalisation := 2.

	if xComp == 0 && yComp == 0 {
		normalisation = 1.
	}

	thetaX := math.Pi * float64(xComp) / float64(width)
	thetaY := math.Pi * float64(yComp) / float64(height)
	cosXs := make([]float64, width)
	for x := range width {
		cosXs[x] = math.Cos(thetaX * float64(x))
	}

	cosYs := make([]float64, height)
	for y := range height {
		cosYs[y] = math.Cos(thetaY * float64(y))
	}

	// Tight inner loops: avoid y*width for every x and reuse row offset
	for y := range height {
		row := y * width
		for x := range width {
			basis := cosXs[x] * cosYs[y]
			p := pixels[row+x]
			result.r += basis * sRGBToLinear[p.r]
			result.g += basis * sRGBToLinear[p.g]
			result.b += basis * sRGBToLinear[p.b]
		}
	}

	scale := normalisation / float64(width*height)
	result.r *= scale
	result.g *= scale
	result.b *= scale
	*fct = result
}

func Decode() {

}

func IsBlurhashValid(blurhash string) bool {
	return false
}

func linearTosRGB(value float64) int {
	if value < 0 {
		return 0
	} else if value < 0.0031308 {
		return int(value*12.92*255 + 0.5)
	} else if value > 1 {
		return 255
	}

	return int((1.055*math.Pow(value, 0.41666666666)-0.055)*255 + 0.5)
}

func encodeDC(value factor) int {
	roundedR := linearTosRGB(value.r)
	roundedG := linearTosRGB(value.g)
	roundedB := linearTosRGB(value.b)
	return (roundedR << 16) + (roundedG << 8) + roundedB
}

func encodeAC(value factor, mx float64) int {
	quantR := max(0, min(18, int(math.Floor(signPow(value.r/mx, 0.5)*9+9.5))))
	quantG := max(0, min(18, int(math.Floor(signPow(value.g/mx, 0.5)*9+9.5))))
	quantB := max(0, min(18, int(math.Floor(signPow(value.b/mx, 0.5)*9+9.5))))
	return quantR*19*19 + quantG*19 + quantB
}

func signPow(val float64, exp float64) float64 {
	return math.Copysign(math.Pow(math.Abs(val), exp), val)
}

func rGBAToPixels(img image.RGBA) (pixels []pixel, width int, height int) {
	size := img.Rect.Max
	pixels = make([]pixel, size.X*size.Y)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			i := img.PixOffset(x, y)
			pixels[y*size.X+x] = pixel{
				r: img.Pix[i],
				g: img.Pix[i+1],
				b: img.Pix[i+2],
			}
		}
	}

	return pixels, size.X, size.Y
}

func yCbCrToPixels(img image.YCbCr) (pixels []pixel, width int, height int) {
	size := img.Rect.Max
	pixels = make([]pixel, size.X*size.Y)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			yi := img.YOffset(x, y)
			ci := img.COffset(x, y)
			r, g, b := color.YCbCrToRGB(img.Y[yi], img.Cb[ci], img.Cr[ci])
			pixels[y*size.X+x] = pixel{
				r: r,
				g: g,
				b: b,
			}
		}
	}

	return pixels, size.X, size.Y
}

func imageToPixels(img image.Image) (pixels []pixel, width int, height int) {
	if img, ok := img.(*image.RGBA); ok {
		return rGBAToPixels(*img)
	}
	if img, ok := img.(*image.YCbCr); ok {
		return yCbCrToPixels(*img)
	}

	size := img.Bounds().Max
	pixels = make([]pixel, size.X*size.Y)
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels[y*size.X+x] = pixel{
				r: uint8(r >> 8),
				g: uint8(g >> 8),
				b: uint8(b >> 8),
			}
		}
	}
	return pixels, size.X, size.Y
}
