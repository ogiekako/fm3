package fm3

import (
	"sort"
	"strconv"
	"strings"
)

type Piece struct {
	Hither bool
	Name   Name // class name or empty
}

type Circum struct {
	Board [10][10]*Piece       // col, row
	Hand map[bool]map[Name]int // hither -> class name -> num
	Hither          bool
	LastIsUchifu    bool
	theHash         int64
}

// String return human readable board.
// TODO(oka): Match the fm's output http://www.abz.jp/~k7ro/report/ki4slct.html
func (c *Circum) String() string {
	res := ""
	res += "後手："+handToStr(c.Hand[false])+"\n"
	var rows []string
	edge := "+---------------------------+"
	rows = append(rows, edge)
	for j := 1; j <= 9; j++ {
		var row []string
		for i := 9; i >= 1; i-- {
			p := c.Board[i][j]
			if p == nil {
				row = append(row, " "+NONE.letter())
			} else if p.Hither {
				row = append(row, " "+p.Name.letter())
			} else {
				row = append(row, "v"+p.Name.letter())
			}
		}
		rows = append(rows, "|" + strings.Join(row, "") + "|")
	}
	rows = append(rows, edge)
	res += strings.Join(rows, "\n")+"\n"
	res += "持駒："+handToStr(c.Hand[true])+"\n"
	return res
}

func handToStr(hand map[Name]int) string {
	var names []int
	for n, _ := range hand {
		names = append(names, int(n))
	}
	sort.Ints(names)
	//	fmt.Println("names", names, hand)

	res := ""
	for _, ni := range names {
		n := Name(ni)
		i := hand[n]
		if i < 0 {
			panic("i < 0")
		}
		if i == 1 {
			res += n.letter()+" "
		} else if i > 1 {
			res += n.letter()+strconv.Itoa(i)+" "
		}
	}
	return res
}

func (c *Circum) clone() *Circum {
	d := NewCircum()
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			d.Board[i][j] = c.Board[i][j]
		}
	}
	for name, num := range c.Hand[true] {
		d.Hand[true][name] = num
	}
	for name, num := range c.Hand[false] {
		d.Hand[false][name] = num
	}
	d.Hither = c.Hither
	d.LastIsUchifu = c.LastIsUchifu
	return d
}

func (c *Circum) hash() int64 {
	if c.theHash != 0 {
		return c.theHash
	}
	const mul = int64(1e9 + 7)
	res := int64(0)
	if c.Hither {
		res += 1
	}
	if c.LastIsUchifu {
		res += 2
	}
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			res *= mul
			if p := c.Board[i][j]; p != nil {
				res += int64(p.Name)*2
				if p.Hither {
					res+= 1
				}
			}
		}
	}
	for n, i := range c.Hand[true] {
		res *= mul
		res +=int64(n)*2
		res += 1
		res *= mul
		res += int64(i)
	}
	for n, i := range c.Hand[false] {
		res *= mul
		res += int64(n)*2
		res *= mul
		res += int64(i)
	}
	c.theHash = res
	return res
}

func NewCircum() *Circum {
	c := &Circum{
		Hand: make(map[bool]map[Name]int),
	}
	c.Hand[true] = make(map[Name]int)
	c.Hand[false] = make(map[Name]int)
	return c
}

// Mated returns if I am mated.
func (c *Circum) Mated() bool {
	return c.Checked() && len(c.Nexts()) == 0
}

// Checked returns if I am checked.
func (c *Circum) Checked() bool {
	me := c.Hither
	you := !me
	// When it is your turn
	c.Hither = you
	defer func() {
		c.Hither = me
	}()
	//	fmt.Println("c",c)
	for _, nc := range c.maybeNexts() {
		// You captured my king
		if nc.Hand[you][OU] > 0 {
			return true
		}
	}
	return false
}

// Illegal returns if this circum is illegal. Illegal circums are
// disregard of check, immovable pieces or uchifuzume.
func (c *Circum) Illegal() bool {
	me := c.Hither
	you := !me
	// disregard of check
	c.Hither = you
	if c.Checked() {
		c.Hither = me
		return true
	}
	c.Hither = me

	// If I am yonder, must be being checked.
	if !me && !c.Checked() {
		return true
	}

	// you have immovable pieces
	if you { // hither
		for col := 1; col <= 9; col++ {
			for row := 1; row <= 2; row++ {
				if p := c.Board[col][row]; p != nil && p.Hither == you {
					if (p.Name == FU || p.Name == KY) && row == 1 {
						return true
					} else if p.Name == KE {
						return true
					}
				}
			}
		}
	} else {// yonder
		for col := 1; col <= 9; col++ {
			for row := 8; row <= 9; row++ {
				if p := c.Board[col][row]; p != nil && p.Hither == you {
					if (p.Name == FU || p.Name == KY) && row == 9 {
						return true
					} else if p.Name == KE {
						return true
					}
				}
			}
		}
	}

	// you made nifu
	for col := 1; col <= 9; col++ {
		found := false
		for row := 1; row <= 9; row++ {
			if p := c.Board[col][row]; p != nil && p.Hither == you && p.Name == FU {
				if found {
					return true
				}
				found = true
			}
		}
	}

	// uchifuzume
	if c.LastIsUchifu && c.Mated() {
		return true
	}
	return false
}

// Nexts returns all next legal circums.
func (c *Circum) Nexts() []*Circum {
	var res []*Circum
	for _, nc := range c.maybeNexts() {
		if !nc.Illegal() {
			res = append(res, nc)
		}
	}
	return res
}

var promotableRow map[bool]map[int]bool = map[bool]map[int]bool {
	true: map[int]bool {
		1: true, 2: true, 3: true,
	},
	false: map[int]bool {
		7: true, 8: true, 9: true,
	},
}

// maybeNexts doesn't consider illegality. May capture king and take it in hand.
func (c *Circum) maybeNexts() []*Circum {
	var res []*Circum
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			p := c.Board[i][j];
			if (p != nil && c.Hither == p.Hither) {
				cl := p.Name.class()
				// leap
				for _, l := range cl.leap {
					if !c.Hither {
						l = V {
							c: -l.c,
							r: -l.r,
						}
					}
					ni := i + l.c
					nj := j + l.r
					if (ni < 1 || 9 < ni || nj < 1 || 9 < nj) {
						continue
					}
					q := c.Board[ni][nj]
					if q == nil || q.Hither != c.Hither {
						res = append(res, move(c, newV(i, j), newV(ni, nj), false))
						if cl.promote != NONE {
							rs := promotableRow[c.Hither]
							if rs[j] || rs[nj] {
								res = append(res, move(c, newV(i, j), newV(ni, nj), true))
							}
						}
					}
				}
				// ride
				for _, r := range cl.ride {
					if !c.Hither {
						r = V {
							c: -r.c,
							r: -r.r,
						}
					}
					for ni, nj := i + r.c, j + r.r; 1 <= ni && ni <= 9 && 1 <= nj && nj <= 9; ni, nj = ni+r.c, nj+r.r {
						q := c.Board[ni][nj]
						if q == nil {
							res = append(res, move(c, newV(i, j), newV(ni, nj), false))
							if cl.promote != NONE {
								rs := promotableRow[c.Hither]
								if rs[j] || rs[nj] {
									res = append(res, move(c, newV(i, j), newV(ni, nj), true))
								}
							}
						} else {
							if q.Hither != c.Hither {
								res = append(res, move(c, newV(i, j), newV(ni, nj), false))
								if cl.promote != NONE {
									rs := promotableRow[c.Hither]
									if rs[j] || rs[nj] {
										res = append(res, move(c, newV(i, j), newV(ni, nj), true))
									}
								}
							}
							break
						}
					}
				}
			}
		}
	}
	for name, num := range c.Hand[c.Hither] {
		if num < 0 {
			panic("num < 0")
		}
		if num == 0 {
			continue
		}
		for i := 1; i <= 9; i++ {
			for j := 1; j <= 9; j++ {
				if c.Board[i][j] == nil {
					d := c.clone()
					d.Hither = !d.Hither
					d.Hand[c.Hither][name]--
					d.Board[i][j] = &Piece {
						Hither: c.Hither,
						Name: name,
					}
					d.LastIsUchifu = name == FU
					res = append(res, d)
				}
			}
		}
	}
	return res
}

func move(c *Circum, from, to V, promote bool) *Circum {
	p := c.Board[from.c][from.r]
	q := c.Board[to.c][to.r]
	me := c.Hither
	you := !me

	d := c.clone()
	d.Hither = you
	d.LastIsUchifu = false
	if q != nil {
		if q.Hither == me {
			panic("Can't capture my piece.")
		}
		if unp := q.Name.class().unpromote; unp != NONE {
			d.Hand[me][unp]++
		} else {
			d.Hand[me][q.Name]++
		}
	}
	if promote {
		d.Board[to.c][to.r] = &Piece {
			Hither: me,
			Name: p.Name.class().promote,
		}
	} else {
		d.Board[to.c][to.r] = p
	}
	d.Board[from.c][from.r] = nil
	return d
}
