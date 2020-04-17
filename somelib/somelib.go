package somelib

import "strings"

const check = "already exists"

func IsExistsErr(b []byte) bool {
	return strings.Contains(string(b), check)
}
