package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

//----------------------------------------
// Knuth's Algorithm X with Dancing Links
//----------------------------------------

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

	fmt.Fprintf(wr, "Number of solutions: %d\n\n", len(d.solutions))
	for _, e := range d.solutions {
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
