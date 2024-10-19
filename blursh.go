package blursh

import (
	"errors"
	"image"
	"image/color"
	"math"
	"strings"
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

	for i := range factors {
		y := i / xComp
		x := i % xComp

		factors[i] = multiplyBasisFunction(x, y, width, height, pixels)
	}

	dc := factors[0]
	ac := factors[1:]
	builder := strings.Builder{}
	builder.Grow(4 + xComp*yComp*2)
	sizeFlag := (xComp - 1) + (yComp-1)*9
	encode83(&builder, sizeFlag, 1)
	maximumValue := 1.

	if len(ac) > 0 {
		actualMaximumValue := 0.
		for _, factor := range ac {
			actualMaximumValue = math.Max(math.Max(math.Max(factor.r, factor.g), factor.b), actualMaximumValue)
		}

		quantisedMaximumValue := math.Max(0, math.Min(82, actualMaximumValue*166-0.5))
		maximumValue = (quantisedMaximumValue + 1) / 166
		encode83(&builder, int(quantisedMaximumValue), 1)
	} else {
		encode83(&builder, 0, 1)
	}

	encode83(&builder, encodeDC(dc), 4)

	for _, factor := range ac {
		encode83(&builder, encodeAC(factor, maximumValue), 2)
	}

	return builder.String(), nil
}

func multiplyBasisFunction(xComp int, yComp int, width int, height int, pixels []pixel) factor {
	result := factor{}
	normalisation := 2.

	if xComp == 0 && yComp == 0 {
		normalisation = 1.
	}

	cosXs := make([]float64, width)

	for x := range cosXs {
		cosXs[x] = math.Cos(math.Pi * float64(xComp) * float64(x) / float64(width))
	}

	for y := 0; y < height; y++ {
		cosY := math.Cos(math.Pi * float64(yComp) * float64(y) / float64(height))

		for x := 0; x < width; x++ {
			basis := cosXs[x] * cosY
			pixel := pixels[y*width+x]
			result.r += basis * sRGBToLinear[pixel.r]
			result.g += basis * sRGBToLinear[pixel.g]
			result.b += basis * sRGBToLinear[pixel.b]
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

func encodeAC(value factor, max float64) int {
	quantR := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(value.r/max, 0.5)*9+9.5)))))
	quantG := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(value.g/max, 0.5)*9+9.5)))))
	quantB := int(math.Floor(math.Max(0, math.Min(18, math.Floor(signPow(value.b/max, 0.5)*9+9.5)))))
	return quantR*19*19 + quantG*19 + quantB
}

func signPow(val float64, exp float64) float64 {
	return math.Copysign(math.Pow(math.Abs(val), exp), val)
}

func imageToPixels(img image.Image) (pixels []pixel, width int, height int) {
	maxPoint := img.Bounds().Max
	pixels = make([]pixel, maxPoint.X*maxPoint.Y)

	switch img := img.(type) {
	case *image.RGBA:
		for y := 0; y < maxPoint.Y; y++ {
			for x := 0; x < maxPoint.X; x++ {
				i := img.PixOffset(x, y)
				s := img.Pix[i : i+3 : i+3]
				pixels[y*maxPoint.X+x] = pixel{
					r: s[0],
					g: s[1],
					b: s[2],
				}
			}
		}
	case *image.YCbCr:
		for y := 0; y < maxPoint.Y; y++ {
			for x := 0; x < maxPoint.X; x++ {
				yi := img.YOffset(x, y)
				ci := img.COffset(x, y)
				r, g, b := color.YCbCrToRGB(img.Y[yi], img.Cb[ci], img.Cr[ci])
				pixels[y*maxPoint.X+x] = pixel{
					r: r,
					g: g,
					b: b,
				}
			}
		}
	default:
		for y := 0; y < maxPoint.Y; y++ {
			for x := 0; x < maxPoint.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				pixels[y*maxPoint.X+x] = pixel{
					r: uint8(r >> 8),
					g: uint8(g >> 8),
					b: uint8(b >> 8),
				}
			}
		}
	}

	return pixels, maxPoint.X, maxPoint.Y
}
