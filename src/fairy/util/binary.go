package util

import "encoding/binary"

func GetUint16(buffer []byte, littleEndian bool) uint16 {
	if littleEndian {
		return binary.LittleEndian.Uint16(buffer)
	} else {
		return binary.BigEndian.Uint16(buffer)
	}
}

func GetUint32(buffer []byte, littleEndian bool) uint32 {
	if littleEndian {
		return binary.LittleEndian.Uint32(buffer)
	} else {
		return binary.BigEndian.Uint32(buffer)
	}
}

func GetUint64(buffer []byte, littleEndian bool) uint64 {
	if littleEndian {
		return binary.LittleEndian.Uint64(buffer)
	} else {
		return binary.BigEndian.Uint64(buffer)
	}
}

func PutUint16(buffer []byte, value uint16, littleEndian bool) {
	if littleEndian {
		binary.LittleEndian.PutUint16(buffer, value)
	} else {
		binary.BigEndian.PutUint16(buffer, value)
	}
}

func PutUint32(buffer []byte, value uint32, littleEndian bool) {
	if littleEndian {
		binary.LittleEndian.PutUint32(buffer, value)
	} else {
		binary.BigEndian.PutUint32(buffer, value)
	}
}

func PutUint64(buffer []byte, value uint64, littleEndian bool) {
	if littleEndian {
		binary.LittleEndian.PutUint64(buffer, value)
	} else {
		binary.BigEndian.PutUint64(buffer, value)
	}
}
