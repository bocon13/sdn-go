package p4info

import "encoding/binary"

func RoundedByte(bitwidth int32) int32 {
	bytes := bitwidth / 8
	if bitwidth%8 > 0 {
		bytes++
	}
	return bytes
}

func Uint8(v uint8) []byte {
	return []byte{v}
}

func Uint16(v uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, v)
	return bytes
}

func Uint32(v uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, v)
	return bytes
}

func Uint64(v uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, v)
	return bytes
}
