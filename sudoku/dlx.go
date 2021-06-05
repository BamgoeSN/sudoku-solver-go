package main

import "math"

type Node struct {
	row, size  int
	col        *Node
	u, d, l, r *Node
}

type DLX struct {
	head      *Node
	column    []*Node
	solutions [][]int
}

func NewDLX(arr [][]bool) *DLX {
	d := new(DLX)
	d.head = new(Node)

	// Initialize column lists
	n := len(arr[0])
	column := make([]*Node, n)
	for i := range column {
		column[i] = new(Node)
	}
	d.head.r, d.head.l = column[0], column[n-1]

	for i := range column {
		column[i].size = 0
		column[i].u, column[i].d = column[i], column[i]
		if i == 0 {
			column[i].l = d.head
		} else {
			column[i].l = column[i-1]
		}
		if i == n-1 {
			column[i].r = d.head
		} else {
			column[i].r = column[i+1]
		}
	}

	// Creating the rest of matrix elements
	for i := 0; i < len(arr); i++ {
		var prev *Node = nil
		for j := 0; j < n; j++ {
			if arr[i][j] {
				curr := new(Node)
				curr.row = i
				curr.col = column[j]
				curr.u, curr.d = column[j].u, column[j]

				if prev != nil {
					curr.l, curr.r = prev, prev.r
					prev.r.l = curr
					prev.r = curr
				} else {
					curr.l = curr
					curr.r = curr
				}

				column[j].u.d = curr
				column[j].u = curr
				column[j].size++
				prev = curr
			}
		}
	}

	d.column = column
	d.solutions = make([][]int, 0)
	return d
}

func (d *DLX) cover(c *Node) {
	c.r.l = c.l
	c.l.r = c.r
	for it := c.d; it != c; it = it.d {
		for jt := it.r; jt != it; jt = jt.r {
			jt.d.u = jt.u
			jt.u.d = jt.d
			jt.col.size--
		}
	}
}

func (d *DLX) uncover(c *Node) {
	for it := c.u; it != c; it = it.u {
		for jt := it.l; jt != it; jt = jt.l {
			jt.d.u = jt
			jt.u.d = jt
			jt.col.size++
		}
	}
	c.r.l = c
	c.l.r = c
}

func (d *DLX) search(k int, solution *[]int) {
	if d.head.r == d.head {
		res := make([]int, len(*solution))
		copy(res, *solution)
		d.solutions = append(d.solutions, res)
		return
	}

	// Linear search for the leftmost column with the least ones
	var ptr *Node = nil
	var low int = math.MaxInt32
	for it := d.head.r; it != d.head; it = it.r {
		if it.size < low {
			low = it.size
			ptr = it
		}
	}
	d.cover(ptr)

	for it := ptr.d; it != ptr; it = it.d {
		// Removing columns with already occupied ones
		*solution = append(*solution, it.row)
		for jt := it.r; jt != it; jt = jt.r {
			d.cover(jt.col)
		}

		// Recursion
		d.search(k+1, solution)
		*solution = (*solution)[:len(*solution)-1]
		// Putting back the removed columns
		for jt := it.l; jt != it; jt = jt.l {
			d.uncover(jt.col)
		}
	}
	d.uncover(ptr)
}

func (d *DLX) Solve() {
	solution := make([]int, 0)
	d.search(0, &solution)
}
