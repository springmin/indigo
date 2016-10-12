package indigo

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mid = func() (uint16, error) {
		return math.MaxUint16, nil
	}
	tc = map[uint64]string{
		0:              "1",
		57:             "z",
		math.MaxUint8:  "5Q",
		math.MaxUint16: "LUv",
		math.MaxUint32: "7YXq9G",
		math.MaxUint64: "jpXCZedGfVQ",
	}
)

func TestNew(t *testing.T) {

	s := Settings{
		StartTime: time.Now(),
		MachineID: mid,
	}

	g := New(s)
	assert.Equal(t, "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz", string(g.base))

	ripple := "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"

	s.Base = []byte(ripple)

	g = New(s)
	assert.Equal(t, ripple, string(g.base))
}

func TestGenerator_NextID(t *testing.T) {

	g := New(Settings{
		StartTime: time.Now(),
		MachineID: mid,
	})

	id1, err := g.NextID()
	require.NoError(t, err)
	assert.NotEmpty(t, id1)

	id2, err := g.NextID()
	require.NoError(t, err)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
}

func TestGenerator_Decompose(t *testing.T) {

	g := New(Settings{
		StartTime: time.Now(),
		MachineID: mid,
	})

	m, err := g.Decompose("KGuFE14P")
	require.NoError(t, err)
	require.NotEmpty(t, m)
	assert.NotEmpty(t, m["id"])

	_, err = g.Decompose("")
	require.Error(t, err)
}

func TestGenerator_NextID_Race(t *testing.T) {

	g := New(Settings{
		StartTime: time.Now(),
		MachineID: mid,
	})

	gs := 2048

	var wg sync.WaitGroup
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			defer wg.Done()
			id, err := g.NextID()
			require.NoError(t, err)
			require.NotEmpty(t, id)
		}()
	}

	wg.Wait()
}

func TestGenerator_NextID_SortIDs(t *testing.T) {

	g := New(Settings{
		StartTime: time.Now(),
		MachineID: mid,
	})

	ids := make([]string, 10)

	var err error
	for i := range ids {
		time.Sleep(10 * time.Millisecond)
		ids[i], err = g.NextID()
		require.NoError(t, err)
	}

	for i := range ids {
		j := rand.Intn(i + 1)
		ids[i], ids[j] = ids[j], ids[i]
	}

	old := make([]string, 10)

	copy(old, ids)
	require.Equal(t, old, ids)

	sort.Strings(ids)
	require.NotEqual(t, old, ids)

	var prev uint64
	for i := range ids {
		m, err := g.Decompose(ids[i])
		require.NoError(t, err)
		require.True(t, prev < m["time"])
		prev = m["time"]
	}
}

func BenchmarkGenerator_NextID(b *testing.B) {

	g := New(Settings{
		StartTime: time.Now(),
		MachineID: mid,
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.NextID()
	}
}
