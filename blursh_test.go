package blursh_test

import (
	"image"
	_ "image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zar4za/blursh"
)

func setUp() (image.Image, error) {
	file, _ := os.Open("gopher.jpg")
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func TestEncode(t *testing.T) {
	expected := []string{
		"LLGSox_N=_n3^+M_R+bcofayWBof", // blurha.sh
		"LLGcAQ_N=_n3^+M_R+bcofayWBof", // blurha.sh
		"LLGcAQ_N=_iv^+MxR+bcoff6WBof", // python package https://github.com/woltapp/blurhash-python
		"LLGcAQ_N=_n3^+MxR+bcoff6WBof", // no visual differences
		"LLGcAQ_N=_n3^+M_R+bcoff6WBof",
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
