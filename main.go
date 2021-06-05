package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

//----------------------------------------
// Knuth's Algorithm X with Dancing Links
//----------------------------------------

type Node struct {
	row, size  int
	col        *Node
	u, d, l, r *Node
}

type DLX struct {
	head      *Node
	column    []*Node
	Solutions [][]int
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
	d.Solutions = make([][]int, 0)
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
		d.Solutions = append(d.Solutions, res)
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

//------
// Main
//------

func main() {
	sc.Split(bufio.ScanWords)
	defer wr.Flush()

	sudoku := make([][]int, 9)
	for i := range sudoku {
		sudoku[i] = make([]int, 9)
		for j := range sudoku[i] {
			sudoku[i][j] = nextInt()
		}
	}

	arr := make([][]bool, 9*9*9)
	for i := range arr {
		arr[i] = make([]bool, 81+27*9)
	}

	for r := range sudoku {
		for c := range sudoku[r] {
			for i := 0; i < 9; i++ {
				arr[81*i+9*r+c][9*r+c] = true
				arr[81*i+9*r+c][81+27*i+r] = true
				arr[81*i+9*r+c][81+27*i+9+c] = true
				arr[81*i+9*r+c][81+27*i+18+(r/3)*3+(c/3)] = true
			}
		}
	}

	d := NewDLX(arr)
	// Remove alreay occupied rows
	for r := range sudoku {
		for c := range sudoku[r] {
			i := sudoku[r][c]
			i--

			if i != -1 {
				targetRow := [4]int{9*r + c, 81 + 27*i + r, 81 + 27*i + 9 + c, 81 + 27*i + 18 + (r/3)*3 + (c / 3)}
				for _, t := range targetRow {
					ptr := d.column[t]
					d.cover(ptr)
				}
			}
		}
	}

	d.Solve()

	fmt.Fprintf(wr, "Number of solutions: %d\n\n", len(d.Solutions))
	for _, e := range d.Solutions {
		ans := make([][]int, 9)
		for i := range ans {
			ans[i] = make([]int, 9)
			copy(ans[i], sudoku[i])
		}
		for _, v := range e {
			c := v % 9
			r := (v / 9) % 9
			i := v / 81
			ans[r][c] = i + 1
		}
		for _, r := range ans {
			fmt.Fprintln(wr, fmt.Sprint(r)[1:18])
		}
		fmt.Fprintln(wr)
	}
}

//---------
// Fast IO
//---------

var fp, _ = os.Open("input.txt")
var sc = bufio.NewScanner(fp)

var ou, _ = os.Create("output.txt")
var wr = bufio.NewWriterSize(ou, 100000)

func nextInt() (res int) {
	sc.Scan()
	text := sc.Text()
	v, _ := strconv.Atoi(text)
	return v
}

func nextInt64() (res int64) {
	sc.Scan()
	text := sc.Text()
	v, _ := strconv.ParseInt(text, 10, 64)
	return v
}
