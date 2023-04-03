package geecache

// A ByteView holds an immutable view of bytes.
// 保存不可变字符视图
type ByteView struct {
	b []byte // 存储真实的缓存值
}

// Len returns the view's length
// 返回视图长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice.
// byteSlice 将数据拷贝作为字节切片返回
// b是只读的，使用byteSlice返回一个拷贝值，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary.
// 将数据作为字符串返回，必要时对其进行复制
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
