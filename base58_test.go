package indigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tc = map[uint64]string{
	0:           "1",
	32:          "Z",
	64:          "27",
	512:         "9q",
	1024:        "Jf",
	16777216:    "2UzHM",
	68719476736: "2ohWHHR",
}

func TestEncodeBase58(t *testing.T) {
	for k, v := range tc {
		assert.Equal(t, v, EncodeBase58(k))
	}
}

func TestDecodeBase58(t *testing.T) {
	for k, v := range tc {
		r, err := DecodeBase58(v)
		require.NoError(t, err)
		assert.Equal(t, k, r)
	}

	_, err := DecodeBase58("0")
	require.Error(t, err)
}