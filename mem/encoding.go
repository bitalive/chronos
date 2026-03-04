package mem

// EncodeUint16 encodes a uint16 value to a byte slice in Big-Endian.
func EncodeUint16(buf []byte, v uint16) {
	_ = buf[1]
	buf[0] = byte(v >> 8)
	buf[1] = byte(v)
}

// EncodeUint32 encodes a uint32 value to a byte slice in Big-Endian.
func EncodeUint32(buf []byte, v uint32) {
	_ = buf[3]
	buf[0] = byte(v >> 24)
	buf[1] = byte(v >> 16)
	buf[2] = byte(v >> 8)
	buf[3] = byte(v)
}

// EncodeInt64 encodes an int64 value to a byte slice in Big-Endian.
func EncodeInt64(buf []byte, v int64) {
	_ = buf[7]
	buf[0] = byte(v >> 56)
	buf[1] = byte(v >> 48)
	buf[2] = byte(v >> 40)
	buf[3] = byte(v >> 32)
	buf[4] = byte(v >> 24)
	buf[5] = byte(v >> 16)
	buf[6] = byte(v >> 8)
	buf[7] = byte(v)
}

// DecodeUint16 decodes a byte slice to a uint16 value in Big-Endian.
func DecodeUint16(buf []byte) uint16 {
	_ = buf[1]
	return uint16(buf[0])<<8 | uint16(buf[1])
}

// DecodeUint32 decodes a byte slice to a uint32 value in Big-Endian.
func DecodeUint32(buf []byte) uint32 {
	_ = buf[3]
	return uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
}

// DecodeInt64 decodes a byte slice to an int64 value in Big-Endian.
func DecodeInt64(buf []byte) int64 {
	_ = buf[7]
	return int64(buf[0])<<56 | int64(buf[1])<<48 | int64(buf[2])<<40 | int64(buf[3])<<32 |
		int64(buf[4])<<24 | int64(buf[5])<<16 | int64(buf[6])<<8 | int64(buf[7])
}

// DecodeUint64 decodes a byte slice to a uint64 value in Big-Endian.
func DecodeUint64(buf []byte) uint64 {
	_ = buf[7]
	return uint64(buf[0])<<56 | uint64(buf[1])<<48 | uint64(buf[2])<<40 | uint64(buf[3])<<32 |
		uint64(buf[4])<<24 | uint64(buf[5])<<16 | uint64(buf[6])<<8 | uint64(buf[7])
}
