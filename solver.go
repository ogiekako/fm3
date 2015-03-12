package fm3

import (
	"errors"
	"fmt"
)

type Solver struct {
	saw map[int64]bool // hash values of seen circums. TODO: Use Circum struct as the key to avoid possible hash collision.
	debug bool
}

func NewSolver() *Solver {
	return &Solver{
		saw: make(map[int64]bool),
	}
}

type Prob struct {
	board *Circum
	step int
}

// solve returns if there is at least one solution.
// returns error too if there is more than one solutions.
func (s *Solver) Solve(prob *Prob) (bool, error) {
	return s.solve(prob.board, prob.step)
}

func (s *Solver) solve(prob *Circum, step int) (bool, error) {
	if s.debug {
		fmt.Printf("Solving\n%d手\n%s", step, prob)
	}
	if !prob.Hither {
		return false, errors.New("turn must be hither.")
	}
	if step%2 != 1 {
		return false, errors.New("step must be odd.")
	}
	var q []*Circum
	q = append(q, prob)
	for i := 0; i < step; i++ {
		if s.debug {
			fmt.Printf("step %d:  #states: %d\n", i, len(q))
		}
		var nq []*Circum
		for _, c := range q {
			if !c.Hither && c.Mated() {
				return true, nil // early mate. TODO: return some error
			}
			for _, nc := range c.Nexts() {
				// TODO: Good branch cut option like FM's -m option.
				if !s.saw[nc.hash()] {
					nq = append(nq, nc)
					s.saw[nc.hash()] = true
				}
			}
		}
		q = nq
	}
	found := false
	for _, c := range q {
		if !c.Hither && c.Mated() {
			if s.debug {
				fmt.Printf("Step: %d\n%s\n", step, c)
			}
			if found {
				return true, errors.New("More than one solution.")
			}
			found = true
		}
	}
	return found, nil
}
