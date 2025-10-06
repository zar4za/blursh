package blursh_test

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/jpeg"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zar4za/blursh"
)

//go:embed gopher.jpg
var gopher []byte

func TestEncode(t *testing.T) {
	img, format, err := image.Decode(bytes.NewReader(gopher))
	require.NoError(t, err)
	require.Equal(t, "jpeg", format)

	actual, err := blursh.Encode(img, 4, 3)

	require.NoError(t, err)
	require.Contains(t, []string{
		"LLGSox_N=_n3^+M_R+bcofayWBof", // blurha.sh
		"LLGcAQ_N=_n3^+M_R+bcofayWBof", // blurha.sh
		"LLGcAQ_N=_iv^+MxR+bcoff6WBof", // python package https://github.com/woltapp/blurhash-python
		"LLGcAQ_N=_n3^+MxR+bcoff6WBof", // no visual differences
		"LLGcAQ_N=_n3^+M_R+bcoff6WBof",
	}, actual)
}

func BenchmarkEncode(b *testing.B) {
	img, format, err := image.Decode(bytes.NewReader(gopher))
	require.NoError(b, err)
	require.Equal(b, "jpeg", format)

	for b.Loop() {
		blursh.Encode(img, 4, 3)
	}
}
