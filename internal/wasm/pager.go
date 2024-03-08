package wasm

type pager struct {
	ChunkSize uint32
	head      *chunk
	last      *chunk
}

func (p *pager) Read(offset, byteCount uint32) []byte {
	p.ensureSize(offset + byteCount)
	res := make([]byte, 0, byteCount)
	chunk, offset := p.getChunkAt(offset)
	if chunk == nil {
		panic("unreachable: pager.ensureSize did not ensure sufficient size")
	}
	res = append(res, chunk.Read(offset, byteCount)...)
	chunk = chunk.next
	for chunk != nil {
		res = append(res, chunk.Read(offset, byteCount)...)
		if offset < byteCount {
			break
		}
		offset -= byteCount
		chunk = chunk.next
	}
	return res
}

func (p *pager) Write(offset uint32, val []byte) {
}

// getChunkAt finds the chunk that contains the given address
// and the offset of the address in the given chunk.
//
// For example, if the chunk size is 4 and the given offset is 11,
// it will return third chunk (range from 8 to 12) and 3 (11 minus 8).
func (p *pager) getChunkAt(offset uint32) (*chunk, uint32) {
	chunk := p.head
	for chunk != nil {
		if offset < p.ChunkSize {
			return chunk, offset
		}
		chunk = chunk.next
		offset -= p.ChunkSize
	}
	return nil, 0
}

// ensureSize expands (if needed) the pager to fit that many bytes.
func (p *pager) ensureSize(size uint32) {
	actChunks := p.chunksCount()
	expChunks := size / p.ChunkSize
	if size%p.ChunkSize > 0 {
		expChunks += 1
	}
	if actChunks >= expChunks {
		return
	}
	for i := actChunks; i < expChunks; i++ {
		p.append(&chunk{raw: make([]byte, p.ChunkSize)})
	}
}

func (p *pager) append(c *chunk) {
	if p.head == nil {
		p.head = c
	} else {
		p.last.next = c
	}
	p.last = c
}

// chunksCount returns how many chunks the pager currently contains.
func (p *pager) chunksCount() uint32 {
	var res uint32 = 0
	chunk := p.head
	for chunk != nil {
		res += 1
		chunk = chunk.next
	}
	return res
}

type chunk struct {
	raw  []byte
	next *chunk
}

func (c *chunk) Read(start, byteCount uint32) []byte {
	end := start + byteCount
	if end > uint32(len(c.raw)) {
		end = uint32(len(c.raw))
	}
	return c.raw[start:end]
}
