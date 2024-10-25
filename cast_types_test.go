package pick

import (
	"testing"
)

func BenchmarkBasicTypesStruct(b *testing.B) {
	for range b.N {
		_ = newBasicTypes()
	}
}
