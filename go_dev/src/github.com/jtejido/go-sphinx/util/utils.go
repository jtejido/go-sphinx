package util

import (
	"encoding/binary"
	"io"
	"strings"
)

func Pad(padding int) string {
	if padding > 0 {
		return strings.Repeat(" ", padding)
	}
	return ""
}

func PadWithMinLength(str string, minLength int) string {
	pad := minLength - len(str)
	if pad > 0 {
		return str + strings.Repeat(" ", pad)
	} else if pad < 0 {
		return str[:minLength]
	}
	return str
}

// readLittleEndianInt reads a little-endian integer from the input stream.
func ReadLittleEndianInt(dataStream io.Reader) (int, error) {
	var val uint32
	var buf [4]byte

	// Read 4 bytes from the dataStream
	_, err := io.ReadFull(dataStream, buf[:])
	if err != nil {
		return 0, err
	}

	// Assemble the little-endian integer from the bytes
	val = uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24

	return int(val), nil
}

// swapInteger byte-swaps the given integer to the other endian.
// That is, if this integer is big-endian, it becomes little-endian, and vice-versa.
func SwapInteger(integer int32) int32 {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(integer))
	return int32(binary.BigEndian.Uint32(bytes))
}

/**
 * If a data point is below 'floor' make it equal to floor.
 *
 * @param data  the data to floor
 * @param floor the floored value
 */
func FloorData(data []float32, floor float32) {
	for i := 0; i < len(data); i++ {
		if data[i] < floor {
			data[i] = floor
		}
	}
}

/**
 * If a data point is non-zero and below 'floor' make it equal to floor
 * (don't floor zero values though).
 *
 * @param data the data to floor
 * @param floor the floored value
 */
func NonZeroFloor(data []float32, floor float32) {
	for i := 0; i < len(data); i++ {
		if data[i] != 0.0 && data[i] < floor {
			data[i] = floor
		}
	}
}

/**
 * Normalize the given data.
 *
 * @param data the data to normalize
 */
func Normalize(data []float32) {
	var sum float32
	for _, val := range data {
		sum += val
	}
	if sum != 0.0 {
		for i := 0; i < len(data); i++ {
			data[i] = data[i] / sum
		}
	}
}
