package cell

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/bits"
	"math/rand"
	"testing"
)

func BenchmarkUint(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		c := BeginCell()
		randomI := uint64(rand.Int63())
		b.StartTimer()
		for j := 0; j < 1000; j++ {
			c.StoreUInt(randomI, bits.Len64(randomI))
		}
		b.StopTimer()
	}
}

func BenchmarkFast(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		c := BeginCell()
		randomI := uint64(rand.Int63())

		b.StartTimer()
		for j := 0; j < 1000; j++ {
			c.StoreUIntFast(randomI, bits.Len64(randomI))
		}
		b.StopTimer()
	}
}

func BenchmarkFastRotate(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		c := BeginCell()
		randomI := uint64(rand.Int63())

		b.StartTimer()
		for j := 0; j < 1000; j++ {
			c.StoreUIntFastRotate(randomI, bits.Len64(randomI))
		}
		b.StopTimer()
	}
}

// func BenchmarkFastFor(b *testing.B) {
// 	b.StopTimer()
// 	for i := 0; i < b.N; i++ {
// 		c := BeginCell()
// 		randomI := uint64(rand.Int63())

// 		b.StartTimer()
// 		for j := 0; j < 1000; j++ {
// 			c.StoreUIntFastFor(randomI, bits.Len64(randomI))
// 		}
// 		b.StopTimer()
// 	}
// }

// func BenchmarkFastSimple(b *testing.B) {
// 	b.StopTimer()
// 	for i := 0; i < b.N; i++ {
// 		c := BeginCell()
// 		randomI := uint64(rand.Int63())

// 		b.StartTimer()
// 		for j := 0; j < 1000; j++ {
// 			c.StoreUIntFastSimple(randomI, bits.Len64(randomI))
// 		}
// 		b.StopTimer()
// 	}
// }

func TestUint(t *testing.T) {
	for i := 0; i < 1000; i++ {
		randomI := uint64(rand.Int63())
		c := BeginCell()
		log.Println("i", randomI)
		err := c.StoreUIntFastRotate(randomI, bits.Len64(randomI))
		if err != nil {
			t.Fatal(err)
			return
		}

		res := c.EndCell()

		u, _ := res.BeginParse().LoadUInt(bits.Len64(randomI))
		if err != nil {
			t.Fatal(err)
			return
		}
		if u != randomI {
			t.Fatal(errors.New("Not parsed"))
			return
		}
	}
}

func TestUintError(t *testing.T) {
	randomI := uint64(4739111663495868)
	c := BeginCell()
	err := c.StoreUIntFastRotate(randomI, bits.Len64(randomI))
	if err != nil {
		t.Fatal(err)
		return
	}

	res := c.EndCell()

	u, _ := res.BeginParse().LoadUInt(bits.Len64(randomI))
	if err != nil {
		t.Fatal(err)
		return
	}
	if u != randomI {
		t.Fatal(errors.New("Not parsed"))
		return
	}
}

func TestUintMax(t *testing.T) {
	randomI := uint64(18446744073709551615)
	c := BeginCell()
	err := c.StoreUIntFastRotate(randomI, bits.Len64(randomI))
	if err != nil {
		t.Fatal(err)
		return
	}

	res := c.EndCell()

	u, _ := res.BeginParse().LoadUInt(bits.Len64(randomI))
	if err != nil {
		t.Fatal(err)
		return
	}
	if u != randomI {
		t.Fatal(errors.New("Not parsed"))
		return
	}
}

func TestCell(t *testing.T) {
	c := BeginCell()

	bs := []byte{11, 22, 33}

	err := c.StoreUInt(1, 1)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = c.StoreSlice(bs, 24)
	if err != nil {
		t.Fatal(err)
		return
	}

	amount := uint64(777)
	c2 := BeginCell().MustStoreCoins(amount).EndCell()

	err = c.StoreRef(c2)
	if err != nil {
		t.Fatal(err)
		return
	}

	u38val := uint64(0xAABBCCF)

	err = c.StoreUInt(u38val, 40)
	if err != nil {
		t.Fatal(err)
		return
	}

	boc := c.EndCell().ToBOC()

	cl, err := FromBOC(boc)
	if err != nil {
		t.Fatal(err)
		return
	}

	lc := cl.BeginParse()

	i, err := lc.LoadUInt(1)
	if err != nil {
		t.Fatal(err)
		return
	}

	if i != 1 {
		t.Fatal("1 bit not eq 1")
		return
	}

	bl, err := lc.LoadSlice(24)
	if err != nil {
		t.Fatal(err)
		return
	}

	if !bytes.Equal(bs, bl) {
		t.Fatal("slices not eq:\n" + hex.EncodeToString(bs) + "\n" + hex.EncodeToString(bl))
		return
	}

	u38, err := lc.LoadUInt(40)
	if err != nil {
		t.Fatal(err)
		return
	}

	if u38 != u38val {
		t.Fatal("uint38 not eq")
		return
	}

	ref, err := lc.LoadRef()
	if err != nil {
		t.Fatal(err)
		return
	}

	amt, err := ref.LoadCoins()
	if err != nil {
		t.Fatal(err)
		return
	}

	if amt != amount {
		t.Fatal("coins ref not eq")
		return
	}
}

func TestCell24(t *testing.T) {
	c := BeginCell()

	bs := []byte{11, 22, 33}

	err := c.StoreSlice(bs, 24)
	if err != nil {
		t.Fatal(err)
		return
	}

	lc := c.EndCell().BeginParse()

	res, err := lc.LoadSlice(24)
	if err != nil {
		t.Fatal(err)
		return
	}

	if !bytes.Equal(bs, res) {
		t.Fatal("slices not eq:\n" + hex.EncodeToString(bs) + "\n" + hex.EncodeToString(res))
		return
	}
}

func TestCell25(t *testing.T) {
	c := BeginCell()

	bs := []byte{11, 22, 33, 0x80}

	err := c.StoreSlice(bs, 25)
	if err != nil {
		t.Fatal(err)
		return
	}

	lc := c.EndCell().BeginParse()

	res, err := lc.LoadSlice(25)
	if err != nil {
		t.Fatal(err)
		return
	}

	if !bytes.Equal(bs, res) {
		t.Fatal("slices not eq:\n" + hex.EncodeToString(bs) + "\n" + hex.EncodeToString(res))
		return
	}
}

func TestCellReadSmall(t *testing.T) {
	c := BeginCell()

	bs := []byte{0b10101010, 0x00, 0x00}

	err := c.StoreSlice(bs, 24)
	if err != nil {
		t.Fatal(err)
		return
	}

	lc := c.EndCell().BeginParse()

	for i := 0; i < 8; i++ {
		res, err := lc.LoadUInt(1)
		if err != nil {
			t.Fatal(err)
			return
		}

		if (res != 1 && i%2 == 0) || (res != 0 && i%2 == 1) {
			t.Fatal("not eq " + fmt.Sprint(i*2))
			return
		}
	}

	res, err := lc.LoadUInt(1)
	if err != nil {
		t.Fatal(err)
		return
	}

	if res != 0 {
		t.Fatal("not 0")
		return
	}
}

func TestCellReadEmpty(t *testing.T) {
	c := BeginCell().EndCell().BeginParse()
	sz, _, err := c.RestBits()
	if err != nil {
		t.Fatal(err)
		return
	}

	if sz != 0 {
		t.Fatal("not 0")
		return
	}
}
