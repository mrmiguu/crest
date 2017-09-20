package rest

import "encoding/binary"
import "strconv"

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

func bytes2bool(b []byte) bool {
	tf, _ := strconv.ParseBool(string(b))
	return tf
}

func bool2bytes(b bool) []byte {
	return strconv.AppendBool(nil, b)
}
