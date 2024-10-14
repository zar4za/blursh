package blursh

import (
	"errors"
	"image"
	"image/color"
	"math"
	"strings"
)

const bytesPerPixel = 4

type factor struct {
	r float64
	g float64
	b float64
}

func Encode(img image.Image, xComp int, yComp int) (string, error) {
	if xComp < 1 || xComp > 9 {
		return "", errors.New("xComp must be in range from 1 to 9")
	}
	if yComp < 1 || yComp > 9 {
		return "", errors.New("yComp must be in range from 1 to 9")
	}

	size := img.Bounds().Max
	// bytesPerRow := bytesPerPixel * width
	factors := make([]factor, xComp*yComp)

	for i := range factors {
		y := i % yComp // row
		x := i / yComp // pix in row
		factors[i] = multiplyBasisFunction(x, y, size.X, size.Y, img)
	}

	dc := factors[0]
	ac := factors[1:]
	builder := strings.Builder{}

	sizeFlag := (xComp - 1) + (yComp-1)*9
	Encode83(&builder, sizeFlag, 1)
	maximumValue := 1.

	if len(ac) > 0 {
		actualMaximumValue := 0.
		for _, factor := range ac {
			actualMaximumValue = math.Max(math.Max(math.Max(factor.r, factor.g), factor.b), actualMaximumValue)
		}

		quantisedMaximumValue := math.Floor(math.Max(0, math.Min(82, math.Floor(actualMaximumValue*166-0.5))))
		maximumValue = (quantisedMaximumValue + 1) / 166
		Encode83(&builder, int(quantisedMaximumValue), 1)
	} else {
		Encode83(&builder, 0, 1)
	}

	Encode83(&builder, encodeDC(dc), 4)

	for _, factor := range ac {
		Encode83(&builder, encodeAC(factor, maximumValue), 2)
	}

	return builder.String(), nil
}

func multiplyBasisFunction(xComp int, yComp int, width int, height int, img image.Image) factor {
	result := factor{}
	normalisation := 2.

	if xComp == 0 && yComp == 0 {
		normalisation = 1.
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			basis := math.Cos(math.Pi*float64(xComp)*float64(x)/float64(width)) * math.Cos(math.Pi*float64(yComp)*float64(y)/float64(height))
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			result.r += basis * sRGBToLinear[rgba.R]
			result.g += basis * sRGBToLinear[rgba.G]
			result.b += basis * sRGBToLinear[rgba.B]
		}
	}

	scale := normalisation / float64(width*height)
	result.r *= scale
	result.g *= scale
	result.b *= scale
	return result
}

func Decode() {

}

func IsBlurhashValid(blurhash string) bool {
	return false
}

func linearTosRGB(value float64) int {
	v := math.Max(0, math.Min(1, value))

	if v <= 0.0031308 {
		return int(math.Trunc(v*12.92*255 + 0.5))
	} else {
		return int(math.Trunc((1.055*math.Pow(v, 1/2.4)-0.055)*255 + 0.5))
	}
}

func encodeDC(value factor) int {
	roundedR := linearTosRGB(value.r)
	roundedG := linearTosRGB(value.g)
	roundedB := linearTosRGB(value.b)
	return (roundedR << 16) + (roundedG << 8) + roundedB
}

func encodeAC(value factor, max float64) int {
	quantR := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(value.r/max, 0.5)*9+9.5)))))
	quantG := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(value.g/max, 0.5)*9+9.5)))))
	quantB := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(value.b/max, 0.5)*9+9.5)))))
	return quantR*19*19 + quantG*19 + quantB
}

func signPow(val float64, exp float64) float64 {
	sign := 1.
	if val < 0 {
		sign = -1.
	}

	return sign * math.Pow(math.Abs(val), exp)
}
