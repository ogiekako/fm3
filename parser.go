package fm3

import (
	"strings"
	"regexp"
	"strconv"
	"errors"
	"log"
	"fmt"

)

type Ptr struct {
	p int
}

// Parse parses the input and parse it as a list of problems.
func Parse(lines []string) []*Prob {
	p := Ptr {
		p: 0,
	}
	var res []*Prob
	for ; ; {
		prob, err := parse(lines, &p);
		if err != nil {
			log.Fatal(err)
		}
		if prob == nil {
			return res
		}
		res = append(res, prob)
	}
}

func parse(lines []string, p *Ptr) (*Prob, error) {
	for ; p.p < len(lines) && !strings.Contains(lines[p.p], "ばか詰"); p.p++ {
	}
	if p.p == len(lines) {
		return nil, nil
	}
	stepS := regexp.MustCompile("\\d+").FindString(lines[p.p])
	p.p++
	step, err := strconv.Atoi(stepS)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("step %d is not an integer.", stepS))
	}
	c := parseCircum(lines, p)
	return &Prob {
		board: c,
		step: step,
	}, nil
}

func parseCircum(lines []string, p *Ptr) *Circum {
	edge := "+---------------------------+"
	for !strings.HasPrefix(lines[p.p], edge) {
		p.p++
	}
	p.p++
	res := NewCircum()
	res.Hither = true
	for j := 1; j <= 9; j, p.p = j+1, p.p+1 {
		row := computeRow(lines[p.p])
		for i := 1; i <= 9; i++ {
			res.Board[i][j] = row[i]
		}
	}
	p.p++
	// TODO: hand
	return res
}

func computeRow(line string) [10]*Piece {
	rs := []rune(line)
	var res [10]*Piece
	for i := 9; i >= 1; i-- {
		p := (9 - i) * 2
		if string(rs[p + 2]) != "・" {
			res[i] = &Piece {
				Hither: string(rs[p + 1]) == " ",
				Name: fromLetter[string(rs[p + 2])],
			}
		}
	}
	return res
}
