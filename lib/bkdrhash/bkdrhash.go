package bkdrhash

import "bytes"

// 计算字节序列hash
func BKDRBytesHash(b []byte) uint32 {
	seed := uint32(131)
	hash := uint32(0)

	for _, v := range b {
		hash = hash*seed + uint32(v)
	}
	return hash
}

// 计算字符串hash
func BKDRHash(s string) uint32 {
	b := bytes.NewBufferString(s).Bytes()
	return BKDRBytesHash(b)
}
