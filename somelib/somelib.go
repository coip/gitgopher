package somelib

import "bytes"

var check = []byte("already exists")

func IsExistsErr(b []byte) bool {
	return bytes.Contains(b, check)
}
