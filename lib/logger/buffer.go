package logger

const _size = 1024

type Buffer struct {
	bs   []byte
	pool Pool
}

func (b *Buffer) Reset() {
	b.bs = b.bs[:0]
}

func (b *Buffer) Free() {
	b.pool.put(b)
}

func (b *Buffer) Copy(data []byte) {
	b.bs = make([]byte, len(data))
	copy(b.bs, data)
}
