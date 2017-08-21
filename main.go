package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

// https://en.wikipedia.org/wiki/Wagner%E2%80%93Fischer_algorithm
func ed(s, t []interface{}) {
	m := len(s)
	n := len(t)

	op := make([][]rune, m)
	for i := 0; i < m; i++ {
		op[i] = make([]rune, n)
	}

	d := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		d[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			if reflect.DeepEqual(s[i-1], t[j-1]) {
				d[i][j] = d[i-1][j-1] // no op required
				op[i-1][j-1] = ' '
			} else {
				del := d[i-1][j] + 1
				add := d[i][j-1] + 1
				rep := d[i-1][j-1] + 1

				if rep <= del {
					if rep <= add {
						d[i][j] = rep
						op[i-1][j-1] = 'R'
					} else {
						// add < sub
						d[i][j] = add
						op[i-1][j-1] = 'A' //
					}
				} else {
					// del < sub
					if add <= del {
						d[i][j] = add
						op[i-1][j-1] = 'A' // d[i][j-1]
					} else {
						// del < add
						d[i][j] = del
						op[i-1][j-1] = 'D' // [i-1][j]
					}
				}
			}

		}
	}

	WriteMatrix(s, t, d, os.Stdout)

}

// WriteMatrix writes a visual representation of the given matrix for the given
// strings to the given writer.
func WriteMatrix(source []interface{}, target []interface{}, matrix [][]int, writer io.Writer) {
	fmt.Fprintf(writer, "    ")
	for _, targetRune := range target {
		fmt.Fprintf(writer, "  %v", targetRune)
	}
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "  %2d", matrix[0][0])
	for j := range target {
		fmt.Fprintf(writer, " %2d", matrix[0][j+1])
	}
	fmt.Fprintf(writer, "\n")
	for i, sourceRune := range source {
		fmt.Fprintf(writer, "%v %2d", sourceRune, matrix[i+1][0])
		for j := range target {
			fmt.Fprintf(writer, " %2d", matrix[i+1][j+1])
		}
		fmt.Fprintf(writer, "\n")
	}
}

func main() {
	fmt.Println("aa")

	s := []interface{}{"t", "a", "m", "a", "l"}
	t := []interface{}{"k", "a", "r", "i", "m"}

	ed(s, t)

}
