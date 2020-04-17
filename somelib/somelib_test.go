package somelib_test

import (
	"testing"

	"github.com/coip/gitgopher/somelib"
)

var (
	pos      = []byte("destination path 'moneypenny' already exists and is not an empty directory.")
	negshort = []byte("somethingshort")
	neglong  = []byte("somethingnooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooootshort")
)

func BenchmarkIsExistsErrPositive(b *testing.B) {
	for n := 0; n < b.N; n++ {
		somelib.IsExistsErr(pos)
	}
}

func BenchmarkisExistsErrNegativeShort(b *testing.B) {
	for n := 0; n < b.N; n++ {
		somelib.IsExistsErr(negshort)
	}
}

func BenchmarkisExistsErrNegativeLong(b *testing.B) {
	for n := 0; n < b.N; n++ {
		somelib.IsExistsErr(neglong)
	}
}
