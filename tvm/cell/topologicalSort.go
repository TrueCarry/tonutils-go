package cell

import (
	"encoding/hex"
)

func topologicalSort(src []*Cell) []*Cell {
	pending := src                    //[]*Cell{src}
	allCells := map[string]*Cell{}    // new Map<string, { cell: Cell, refs: string[] }>();
	notPermCells := map[string]bool{} //new Set<string>();
	sorted := []string{}

	for len(pending) > 0 {
		cells := append([]*Cell{}, pending...)
		pending = []*Cell{}

		for _, cell := range cells {
			hash := hex.EncodeToString(cell.Hash())
			if _, ok := allCells[hash]; ok {
				continue
			}

			notPermCells[hash] = true
			allCells[hash] = cell

			// for _, r := range cell.refs {
			pending = append(pending, cell.refs...)
			// }
		}
	}

	tempMark := map[string]bool{}
	var visit func(hash string)
	visit = func(hash string) {
		if !notPermCells[hash] {
			return
		}

		if tempMark[hash] {
			panic("Not a DAG")
		}

		tempMark[hash] = true

		for _, c := range allCells[hash].refs {
			visit(hex.EncodeToString(c.Hash()))
		}

		sorted = append([]string{hash}, sorted...)
		delete(tempMark, hash)
		delete(notPermCells, hash)
	}

	for len(notPermCells) > 0 {
		for k := range notPermCells {
			visit(k)
			break
		}
	}

	indexes := map[string]int{}
	for i := 0; i < len(sorted); i++ {
		indexes[sorted[i]] = i
	}

	result := []*Cell{}
	for _, ent := range sorted {
		rrr := allCells[ent]
		// rrr.index = indexes[hex.EncodeToString(rrr.Hash())]
		result = append(result, rrr) //.push({ cell: rrr.cell, refs: rrr.refs.map((v) => indexes.get(v)!) });
	}

	// we need to do it this way because we can have same cells but 2 diff object pointers
	var indexSetter func(node *Cell)
	indexSetter = func(node *Cell) {
		node.index = indexes[hex.EncodeToString(node.Hash())]
		for _, ref := range node.refs {
			indexSetter(ref)
		}
	}

	for _, root := range result {
		indexSetter(root)
	}

	return result
}
