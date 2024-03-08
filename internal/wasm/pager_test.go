package wasm

import (
	"testing"

	"github.com/tetratelabs/wazero/internal/testing/require"
)

func TestPagerReadAllocation(t *testing.T) {
	p := pager{ChunkSize: 4}
	require.Equal(t, p.chunksCount(), uint32(0))
	var b []byte

	// allocate first chunk
	b = p.Read(0, 2)
	require.Equal(t, p.chunksCount(), uint32(1))
	require.Equal(t, len(b), 2)

	// don't allocate when reading the same range
	b = p.Read(0, 2)
	require.Equal(t, p.chunksCount(), uint32(1))
	require.Equal(t, len(b), 2)

	// don't allocate when reading range from existing chunk
	b = p.Read(1, 2)
	require.Equal(t, p.chunksCount(), uint32(1))
	require.Equal(t, len(b), 2)

	// allocate when offset is in the next chunk
	b = p.Read(5, 2)
	require.Equal(t, p.chunksCount(), uint32(2))
	require.Equal(t, len(b), 2)

	// allocate chunk when offset is at the start of the new chunk
	b = p.Read(7, 2)
	require.Equal(t, p.chunksCount(), uint32(3))
	require.Equal(t, len(b), 2)

	// allocate when the end of the requested range
	// is outside of the last allocated chunk
	b = p.Read(8, 5)
	require.Equal(t, p.chunksCount(), uint32(4))
	require.Equal(t, len(b), 5)

	// allocate as many chunks as needed when reading far ahead
	b = p.Read(72, 2)
	require.Equal(t, p.chunksCount(), uint32(19))
	require.Equal(t, len(b), 2)

}
