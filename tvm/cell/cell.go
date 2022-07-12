package cell

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"strings"
	"sync"
)

type Cell struct {
	special bool
	level   byte
	bitsSz  int
	index   int
	data    []byte

	refs []*Cell
}

func (c *Cell) BeginParse() *LoadCell {
	// copy data
	data := append([]byte{}, c.data...)

	refs := make([]*LoadCell, len(c.refs))
	for i, ref := range c.refs {
		refs[i] = ref.BeginParse()
	}

	return &LoadCell{
		special: c.special,
		level:   c.level,
		bitsSz:  c.bitsSz,
		data:    data,
		refs:    refs,
	}
}

func (c *Cell) ToBuilder() *Builder {
	// copy data
	data := append([]byte{}, c.data...)

	return &Builder{
		bitsSz: c.bitsSz,
		data:   data,
		refs:   c.refs,
	}
}

func (c *Cell) BitsSize() int {
	return c.bitsSz
}

func (c *Cell) RefsNum() int {
	return len(c.refs)
}

func (c *Cell) Dump() string {
	return c.dump(0, false)
}

func (c *Cell) DumpBits() string {
	return c.dump(0, true)
}

func (c *Cell) dump(deep int, bin bool) string {
	sz, data, _ := c.BeginParse().RestBits()

	var val string
	if bin {
		for _, n := range data {
			val += fmt.Sprintf("%08b", n)
		}
		if sz%8 != 0 {
			val = val[:len(val)-(8-(sz%8))]
		}
	} else {
		val = hex.EncodeToString(data)
	}

	str := strings.Repeat("  ", deep) + fmt.Sprint(sz) + "[" + val + "]"
	if len(c.refs) > 0 {
		str += " -> {"
		for i, ref := range c.refs {
			str += "\n" + ref.dump(deep+1, bin)
			if i == len(c.refs)-1 {
				str += "\n"
			} else {
				str += ","
			}
		}
		str += strings.Repeat("  ", deep)
		return str + "}"
	}
	return str
}

var hash256 = sha256.New()
var hashMutex = sync.Mutex{}

func (c *Cell) Hash() []byte {
	serialized := c.serialize(-1, true)

	hashMutex.Lock()
	hash256.Write(serialized)
	sum := hash256.Sum(nil)
	hash256.Reset()
	hashMutex.Unlock()

	return sum
}

func (c *Cell) InternalHash() uint32 {
	h := fnv.New32a()
	h.Write(c.serialize(-1, true))
	return h.Sum32()
}

func (c *Cell) Sign(key ed25519.PrivateKey) []byte {
	return ed25519.Sign(key, c.Hash())
}
