package squads

import (
	"encoding/binary"
)

func toU32Bytes(num uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, num)
	return buf
}

func toU64Bytes(num uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, num)
	return buf
}

// convertToUint8Slice converts []int to []uint8
func convertToUint8Slice(ints []uint16) []uint8 {
	result := make([]uint8, len(ints))
	for i, v := range ints {
		result[i] = uint8(v)
	}
	return result
}
