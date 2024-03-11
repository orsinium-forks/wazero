package wasm

import (
	"testing"

	"github.com/tetratelabs/wazero/internal/testing/require"
)

func TestPagerWrite(t *testing.T) {
	p := pager{ChunkSize: 4}
	require.Equal(t, p.chunksCount(), uint32(0))
	written := p.Write(7, []byte{2, 3, 4, 5})
	require.Equal(t, written, uint32(4))

	b := make([]byte, 9)
	read := p.Read(6, b)
	_ = read
	// require.Equal(t, read, uint32(9))
	require.Equal(t, b, []byte{0, 2, 3, 4, 5, 0, 0, 0, 0})
}
