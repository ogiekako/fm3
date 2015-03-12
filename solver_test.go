package fm3

import (
	"testing"
	"fmt"
)

type onBoard struct {
	hither bool
	c      int
	r      int
	name   Name
}

func on(c, r int, hither bool, name Name) *onBoard {
	return &onBoard{
		hither: hither,
		c:c,
		r:r,
		name:name,
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		descr string
		board []*onBoard
		hand map[bool]map[Name]int
		step  int
		want  bool
		debug bool
	}{{
		descr: "move only",
		board: []*onBoard{
			on(1, 9, false, OU),
			on(1, 8, true, GI),
			on(2, 8, true, KI),
		},
		step: 1,
		want: true,
	}, {
		descr: "I'm checked",
		board: []*onBoard{
			on(1, 9, false, OU),
			on(1, 8, true, FU),
			on(2, 8, true, KI),
			on(9, 9, true, OU),
			on(9, 8, false, KI),
		},
		step: 1,
		want: false,
	}, {
		descr: "promote",
		board: []*onBoard{
			on(1, 3, true, GI),
			on(1, 4, false, OU),
			on(1, 5, false, FU),
			on(2, 5, true, FU),
		},
		step: 1,
		want: true,
	} , {
		descr: "promotion is only possible on oponent's area.",
		board: []*onBoard{
			on(1, 4, true, GI),
			on(1, 5, false, OU),
			on(1, 6, false, FU),
			on(2, 6, true, FU),
		},
		step: 1,
		want: false,
	}, {
		descr: "mate in three",
		board: []*onBoard{
			on(2, 3, true, KI),
			on(1, 4, false, OU),
			on(1, 6, false, FU),
			on(2, 6, true, FU),
		},
		step: 3,
		want: true,
	}, {
		descr: "rider",
		board: []*onBoard{
			on(1, 9, false, OU),
			on(2, 3, true, HI),
			on(2, 9, true, FU),
		},
		step: 1,
		want: true,
	}, {
		descr: "must check",
		board: []*onBoard{
			on(1, 3, false, OU),
			on(2, 5, true, FU),
			on(1, 5, false, FU),
			on(4, 4, true, KI),
		},
		step: 3,
		want:false,
	}, {
		descr: "drop",
		board: []*onBoard {
			on(1, 1, false, OU),
			on(1, 3, true, FU),
		},
		hand: map[bool]map[Name]int{
			true: {
				KI: 1,
			},
		},
		step: 1,
		want: true,
	}, {
		descr: "capture and drop",
		board: []*onBoard{
			on(1, 1, false, OU),
			on(5, 5, false, KA),
			on(9, 5, true, RY),
			on(3, 4, true, KE),
			on(9, 1, false, KI),
		},
		step: 3,
		want: true,
	}, {
		descr: "immovable piece",
		board: []*onBoard {
			on(1, 9, false, OU),
			on(3, 1, true, HI),
			on(1, 8, false, FU),
		},
		hand: map[bool]map[Name]int{
			false: {
				FU: 18,
				KY: 4,
				KE: 4,
			},
		},
		step: 1,
		want: true,
	}, {
		descr: "nifu",
		board: []*onBoard {
			on(5, 2, false, OU),
			on(5, 9, true, FU),
		},
		hand: map[bool]map[Name]int {
			true: {
				FU: 1,
				KI: 1,
			},
		},
		step: 3,
		want: false,
	}, {
		descr: "uchifuzume",
		board: []*onBoard{
			on(1, 1, false, OU),
			on(1, 3, true, GI),
			on(2, 1, false, FU),
		},
		hand: map[bool]map[Name]int {
			true: {
				FU: 1,
			},
		},
		step: 1,
		want: false,
	}, {
		descr: "tsukifuzume",
		board: []*onBoard {
			on(1, 4, false, OU),
			on(2, 2, true, RY),
			on(1, 6, true, FU),
			on(1, 7, true, KY),
		},
		step: 1,
		want: true,
	}, {
		// TODO: this test is too slow.
		descr: "http://www.abz.jp/~k7ro/overflow/hr19.htm",
		board: []*onBoard {
			on(1, 1, false, OU),
			on(1, 2, false, GI),
			on(2, 1, false, GI),
			on(2, 2, false, FU),
			on(3, 1, false, FU),
		},
		hand: map[bool]map[Name]int {
			true: {
				KE: 2,
				FU: 1,
			},
		},
		step: 17,
		want: true,
		debug: true,
	},
	}

	solver := NewSolver()
	for _, tt := range tests {
		if tt.debug {
			fmt.Printf("---------- %q ----------\n", tt.descr)
		}
		prob := NewCircum()
		prob.Hither = true
		for _, b := range tt.board {
			prob.Board[b.c][b.r] = &Piece {
				Hither: b.hither,
				Name: b.name,
			}
		}
		for _, hither := range []bool{true, false} {
			for name, num := range tt.hand[hither] {
				prob.Hand[hither][name] = num
			}
		}
		solver.debug = tt.debug
		got, err := solver.solve(prob, tt.step)
		if err != nil {
			t.Errorf("%q: %v", tt.descr, err)
		}
		if got != tt.want {
			t.Errorf("%q: got: %v  want: %v", tt.descr, got, tt.want)
		}
	}
}
