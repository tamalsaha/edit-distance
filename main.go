package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	lev "github.com/texttheater/golang-levenshtein/levenshtein"
)

// https://en.wikipedia.org/wiki/Wagner%E2%80%93Fischer_algorithm
func ed(s, t []interface{}) {
	m := len(s)
	n := len(t)

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
			} else {
				del := d[i-1][j] + 1
				add := d[i][j-1] + 1
				rep := d[i-1][j-1] + 1
				d[i][j] = min(rep, min(add, del))
			}
		}
	}

	WriteMatrix(s, t, d, os.Stdout)
	es := backtrace(m, n, d)
	for _, op := range es {
		fmt.Print(op.String(), "|")
	}
}

type EditOperation int

const (
	Ins = iota
	Del
	Sub
	Match
)

type EditScript []EditOperation

type MatchFunction func(rune, rune) bool

func (operation EditOperation) String() string {
	if operation == Match {
		return "M"
	} else if operation == Ins {
		return "A"
	} else if operation == Sub {
		return "R"
	}
	return "D"
}

func backtrace(i int, j int, matrix [][]int) EditScript {
	if i > 0 && matrix[i-1][j]+1 == matrix[i][j] {
		return append(backtrace(i-1, j, matrix), Del)
	}
	if j > 0 && matrix[i][j-1]+1 == matrix[i][j] {
		return append(backtrace(i, j-1, matrix), Ins)
	}
	if i > 0 && j > 0 && matrix[i-1][j-1]+1 == matrix[i][j] {
		return append(backtrace(i-1, j-1, matrix), Sub)
	}
	if i > 0 && j > 0 && matrix[i-1][j-1] == matrix[i][j] {
		return append(backtrace(i-1, j-1, matrix), Match)
	}
	return []EditOperation{}
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

func min(x int, y int) int {
	if y < x {
		return y
	}
	return x
}

func max(a int, b int) int {
	if b > a {
		return b
	}
	return a
}

func main() {
	s := []interface{}{"A", "B", "C"}
	t := []interface{}{"D"}

	ed(s, t)
	fmt.Println("")

	es := lev.EditScriptForStrings([]rune("ABC"), []rune("D"), lev.Options{
		InsCost: 1,
		DelCost: 1,
		SubCost: 1,
		Matches: func(sourceCharacter rune, targetCharacter rune) bool {
			return sourceCharacter == targetCharacter
		},
	})
	for _, op := range es {
		fmt.Print(op.String(), "|")
	}
}
