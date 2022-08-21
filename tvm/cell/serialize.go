package cell

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"log"
	"math"
)

var ErrTooBigValue = errors.New("too big value")
var ErrNegative = errors.New("value should be non negative")
var ErrRefCannotBeNil = errors.New("ref cannot be nil")
var ErrSmallSlice = errors.New("too small slice for this size")
var ErrTooBigSize = errors.New("too big size")
var ErrTooMuchRefs = errors.New("too much refs")
var ErrNotFit1023 = errors.New("cell data size should fit into 1023 bits")
var ErrNoMoreRefs = errors.New("no more refs exists")
var ErrAddressTypeNotSupported = errors.New("address type is not supported")

func (c *Cell) ToBOC() []byte {
	return c.ToBOCWithFlags(true)
}

func (c *Cell) ToBOCWithFlags(withCRC bool) []byte {
	// recursively go through cells, build hash index and store unique in slice
	orderCells := flattenIndex([]*Cell{c}) // topologicalSort([]*Cell{c})
	// flattenIndex([]*Cell{c})

	// bytes needed to store num of cells
	cellSizeBits := math.Log2(float64(len(orderCells)) + 1)
	cellSizeBytes := byte(math.Ceil(cellSizeBits / 8))

	var payload []byte
	for i := 0; i < len(orderCells); i++ {
		// serialize each cell
		payload = append(payload, orderCells[i].serialize(uint(cellSizeBytes), false)...)
	}

	// bytes needed to store len of payload
	sizeBits := math.Log2(float64(len(payload)) + 1)
	sizeBytes := byte(math.Ceil(sizeBits / 8))

	// has_idx 1bit, hash_crc32 1bit,  has_cache_bits 1bit, flags 2bit, size_bytes 3 bit
	flags := byte(0b0_0_0_00_000)
	if withCRC {
		flags |= 0b0_1_0_00_000
	}

	flags |= cellSizeBytes

	var data []byte

	data = append(data, bocMagic...)
	data = append(data, flags)

	// bytes needed to store size
	data = append(data, sizeBytes)

	// cells num
	data = append(data, dynamicIntBytes(uint64(calcCells(c)), uint(cellSizeBytes))...)

	// roots num (only 1 supported for now)
	data = append(data, dynamicIntBytes(1, uint(cellSizeBytes))...)

	// complete BOCs = 0
	data = append(data, dynamicIntBytes(0, uint(cellSizeBytes))...)

	// len of data
	data = append(data, dynamicIntBytes(uint64(len(payload)), uint(sizeBytes))...)

	// root should have index 0
	data = append(data, dynamicIntBytes(0, uint(cellSizeBytes))...)
	data = append(data, payload...)

	if withCRC {
		checksum := make([]byte, 4)
		binary.LittleEndian.PutUint32(checksum, crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli)))

		data = append(data, checksum...)
	}

	return data
}

func calcCells(cell *Cell) int {
	m := map[string]*Cell{}
	// calc unique cells
	uniqCells(m, cell)

	return len(m)
}

func uniqCells(m map[string]*Cell, cell *Cell) {
	m[hex.EncodeToString(cell.Hash())] = cell

	for _, ref := range cell.refs {
		uniqCells(m, ref)
	}
}

func flattenIndex(roots []*Cell) []*Cell {
	var indexed []*Cell
	var offset int
	hashIndex := map[string]int{}

	var doIndex func([]*Cell)
	var next [][]*Cell

	doIndex = func(cells []*Cell) {
		for _, c := range cells {
			h := hex.EncodeToString(c.Hash())

			id, ok := hashIndex[h]
			if !ok {
				id = offset
				offset++

				hashIndex[h] = id

				indexed = append(indexed, c)
				if len(c.refs) > 0 {
					next = append(next, c.refs)
				}
			} else { // if we already have such cell, we need to move it forward in order.
				log.Println("Index move", id)

				// before := ""
				// for _, c := range indexed {
				// 	before += fmt.Sprintf("%v,", c.BeginParse().MustLoadUInt(8))
				// }

				next = append(next, indexed[id].refs)
				// move to end
				indexed = append(indexed, indexed[id])

				// remove from old position
				indexed = append(indexed[:id], indexed[id+1:]...)

				// after := ""
				// for _, c := range indexed {
				// 	after += fmt.Sprintf("%v,", c.BeginParse().MustLoadUInt(8))
				// }

				// log.Println("before, after", before, after)

				// reindex
				for i, ci := range indexed {
					// TODO: optimize
					i := i
					th := hex.EncodeToString(ci.Hash())

					hashIndex[th] = i
				}

			}
		}

		// for _, n := range next {
		// 	doIndex(n)
		// }

		// return ordered cells to write to boc
		// return indexed
	}
	// doIndex(roots)
	next = append(next, roots)

	i := 0
	for len(next) > i {
		n := next[i]
		i++
		// for _, k := range next {
		doIndex(n)
		// break
		// }
	}

	// after := ""
	// for _, c := range indexed {
	// 	after += fmt.Sprintf("%v,", c.BeginParse().MustLoadUInt(8))
	// }
	// log.Println("after index", after)

	for i, ci := range indexed {
		// TODO: optimize
		i := i
		th := hex.EncodeToString(ci.Hash())

		hashIndex[th] = i
	}

	// after = ""
	// for _, c := range indexed {
	// 	after += fmt.Sprintf("%v,", c.BeginParse().MustLoadUInt(8))
	// }
	// log.Println("after re index", after)

	// we need to do it this way because we can have same cells but 2 diff object pointers
	var indexSetter func(node *Cell)
	indexSetter = func(node *Cell) {
		node.index = hashIndex[hex.EncodeToString(node.Hash())]
		// log.Println("Node hash", node.Hash(), hashIndex[hex.EncodeToString(node.Hash())], node.BeginParse().MustLoadUInt(8))
		for _, ref := range node.refs {
			indexSetter(ref)
		}
	}

	// after = ""
	// for _, c := range indexed {
	// 	after += fmt.Sprintf("%v,", c.BeginParse().MustLoadUInt(8))
	// }
	// log.Println("after setter", after)

	for _, root := range indexed {
		indexSetter(root)
	}

	return indexed
}

func (c *Cell) serialize(refIndexSzBytes uint, isHash bool) []byte {
	// copy
	payload := append([]byte{}, c.BeginParse().MustLoadSlice(c.bitsSz)...)

	unusedBits := 8 - (c.bitsSz % 8)
	if unusedBits != 8 {
		// we need to set bit at the end if not whole byte was used
		payload[len(payload)-1] += 1 << (unusedBits - 1)
	}

	data := append(c.descriptors(), payload...)

	if !isHash {
		for _, ref := range c.refs {
			data = append(data, dynamicIntBytes(uint64(ref.index), refIndexSzBytes)...)
		}
	} else {
		for _, ref := range c.refs {
			data = append(data, make([]byte, 2)...)
			binary.BigEndian.PutUint16(data[len(data)-2:], uint16(ref.maxDepth(0)))
		}
		for _, ref := range c.refs {
			data = append(data, ref.Hash()...)
		}
	}

	return data
}

// calc how deep is the cell (how long children tree)
func (c *Cell) maxDepth(start int) int {
	d := start
	for _, cc := range c.refs {
		if x := cc.maxDepth(start + 1); x > d {
			d = x
		}
	}
	return d
}

func (c *Cell) descriptors() []byte {
	ceilBytes := c.bitsSz / 8
	if c.bitsSz%8 != 0 {
		ceilBytes++
	}

	// calc size
	ln := ceilBytes + c.bitsSz/8

	specBit := byte(0)
	if c.special {
		specBit = 8
	}

	return []byte{byte(len(c.refs)) + specBit + c.level*32, byte(ln)}
}

func dynamicIntBytes(val uint64, sz uint) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, val)

	return data[8-sz:]
}
