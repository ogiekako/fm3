package fm3

import (
	"testing"
	"io/ioutil"
	"strings"
	"log"
)

func BenchmarkSolve(b *testing.B) {
	dat, err := ioutil.ReadFile("bench.txt")
	if err != nil {
		b.Error(err)
	}
	lines := strings.Split(string(dat), "\n")
	ps := Parse(lines)

	solver := NewSolver()
	solver.debug = true
	for i, prob := range ps {
		log.Printf("Problem %d...\n", i)
		ok, err := solver.solve(prob.board, prob.step)
		if err != nil {
			b.Error(err)
		}
		if !ok {
			b.Error("Problem %d is not OK", i)
		}
	}
}
