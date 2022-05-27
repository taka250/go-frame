package kabucache

type ByteView struct {
	b []byte
}

//Len 返回view的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice返回一份copy
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

//String
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
