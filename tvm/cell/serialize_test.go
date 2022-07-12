package cell

import (
	"math/bits"
	"math/rand"
	"testing"
)

func BenchmarkSerialize(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		c := BeginCell()

		for j := 0; j < 5; j++ {
			randomI := uint64(rand.Int63())
			c.StoreUInt(randomI, bits.Len64(randomI))
		}

		for j := 0; j < 5; j++ {
			dCell := BeginCell()
			for k := 0; k < 3; k++ {
				randomI := uint64(rand.Int63())
				dCell.StoreUInt(randomI, bits.Len64(randomI))
			}

			c.MustStoreRef(dCell.EndCell())
		}

		cell := c.EndCell()

		// for i
		b.StartTimer()
		for j := 0; j < 1000; j++ {

			cell.ToBOCWithFlags(true)

		}
		b.StopTimer()
	}
}
