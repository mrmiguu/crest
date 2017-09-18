package rest

import "encoding/binary"

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func itob(i int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
