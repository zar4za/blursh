package blursh_test

import (
	"blursh"
	"image"
	_ "image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	expected := "LLGcAQ_N=_n3^+M_R+bcofayWBof"

	file, err := os.Open("gopher.jpg")

	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := blursh.Encode(img, 4, 3)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
