package wasm

type Page struct {
	ChunkSize uint32
	head      *chunk
	last      *chunk
}

func (p *Page) Read(offset uint32, buffer []byte) uint32 {
	byteCount := uint32(len(buffer))
	chunk, offset := p.getChunkAt(offset)
	if chunk == nil {
		// start is out of range, just return zeros
		return 0
	}
	read := chunk.Read(offset, buffer)
	buffer = buffer[read:]
	chunk = chunk.next
	total := read
	for chunk != nil {
		read = chunk.Read(0, buffer)
		buffer = buffer[read:]
		total += read
		if byteCount < p.ChunkSize {
			break
		}
		byteCount -= p.ChunkSize
		chunk = chunk.next
	}
	return total
}

func (p *Page) Write(offset uint32, val []byte) uint32 {
	byteCount := uint32(len(val))
	p.ensureSize(offset + byteCount)
	chunk, offset := p.getChunkAt(offset)
	if chunk == nil {
		panic("unreachable: insufficient page size ensured")
	}
	written := chunk.Write(offset, val)
	val = val[written:]
	chunk = chunk.next
	total := written
	for chunk != nil {
		written = chunk.Write(0, val)
		val = val[written:]
		total += written
		if len(val) == 0 {
			break
		}
		chunk = chunk.next
	}
	return total
}

// getChunkAt finds the chunk that contains the given address
// and the offset of the address in the given chunk.
//
// For example, if the chunk size is 4 and the given offset is 11,
// it will return third chunk (range from 8 to 12) and 3 (11 minus 8).
func (p *Page) getChunkAt(offset uint32) (*chunk, uint32) {
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
func (p *Page) ensureSize(size uint32) {
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

func (p *Page) append(c *chunk) {
	if p.head == nil {
		p.head = c
	} else {
		p.last.next = c
	}
	p.last = c
}

// chunksCount returns how many chunks the pager currently contains.
func (p *Page) chunksCount() uint32 {
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

func (c *chunk) Read(start uint32, buffer []byte) uint32 {
	return uint32(copy(buffer, c.raw[start:]))
}

func (c *chunk) Write(start uint32, val []byte) uint32 {
	return uint32(copy(c.raw[start:], val))
}
