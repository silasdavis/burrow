package slurp

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestGzip(t *testing.T) {
	bsIn := []byte("I am a silly frog I am a silly frog I am a silly frog")
	bsComp, err := Gzip(bsIn)
	require.NoError(t, err)
	bsOut, err := Gunzip(bsComp)
	require.NoError(t, err)
	require.Equal(t, bsIn, bsOut)
}
