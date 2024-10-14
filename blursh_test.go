package blursh_test

import (
	"blursh"
	"image"
	_ "image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setUp() (image.Image, error) {
	file, err := os.Open("gopher.jpg")
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func TestEncode(t *testing.T) {
	expected := []string{
		"LLGSox_N=_n3^+M_R+bcofayWBof",
		"LLGcAQ_N=_n3^+M_R+bcofayWBof",
	}

	img, err := setUp()

	if err != nil {
		t.Fatal(err)
	}

	actual, err := blursh.Encode(img, 4, 3)

	assert.Nil(t, err)
	assert.Contains(t, expected, actual)
}

func BenchmarkEncode(b *testing.B) {
	img, err := setUp()

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		blursh.Encode(img, 4, 3)
	}
}
