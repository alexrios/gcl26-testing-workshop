package fuzzing

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

const formatVersion byte = 1

type Entry struct {
	Key   string
	Value string
}

func Encode(e Entry) ([]byte, error) {
	if len(e.Key) > math.MaxUint16 {
		return nil, errors.New("key too large")
	}

	data := make([]byte, 3+len(e.Key)+len(e.Value))
	data[0] = formatVersion
	binary.BigEndian.PutUint16(data[1:3], uint16(len(e.Key)))
	copy(data[3:], e.Key)
	copy(data[3+len(e.Key):], e.Value)
	return data, nil
}

func Decode(data []byte) (Entry, error) {
	if len(data) < 3 {
		return Entry{}, errors.New("header too short")
	}
	if data[0] != formatVersion {
		return Entry{}, fmt.Errorf("unsupported version %d", data[0])
	}

	keyLen := int(binary.BigEndian.Uint16(data[1:3]))
	if keyLen > len(data)-3 {
		return Entry{}, errors.New("truncated key")
	}

	return Entry{
		Key:   string(data[3 : 3+keyLen]),
		Value: string(data[3+keyLen:]),
	}, nil
}
