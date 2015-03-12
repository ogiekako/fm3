package fm3

type V struct {
	c int // delta column
	r int // delta row
}

func newV(c, r int) V {
	return V{
		c: c,
		r: r,
	}
}

type Name int

const (
	NONE = Name(0)
	FU   = Name(1)
	TO   = Name(2)
	KY   = Name(3)
	NKY  = Name(4)
	KE   = Name(5)
	NKE  = Name(6)
	GI   = Name(7)
	NGI  = Name(8)
	KI   = Name(9)
	KA   = Name(10)
	UM   = Name(11)
	HI   = Name(12)
	RY   = Name(13)
	OU   = Name(14)
)

type Class struct {
	name      Name
	leap      []V // hither's point of view ex. FU: (0, -1)
	ride      []V
	king      bool
	unpromote Name // unpromoted name or empty
	promote   Name // promoted name or empty
}

var kinLeap = []V{
	newV(1, -1), newV(0, -1), newV(-1, -1),
	newV(1, 0), /*         */ newV(-1, 0),
	/*        */ newV(0, 1),
}

var toLetter = map[Name]string {
	NONE: "・",
	FU  : "歩",
	TO  : "と",
	KY  : "香",
	NKY : "杏",
	KE  : "桂",
	NKE : "圭",
	GI  : "銀",
	NGI : "全",
	KI  : "金",
	KA  : "角",
	UM  : "馬",
	HI  : "飛",
	RY  : "竜",
	OU  : "玉",
}

var fromLetter = reverse((map[Name]string)(toLetter))

func reverse(map[Name]string) map[string]Name {
	res := make(map[string]Name)
	for n, l := range toLetter {
		res[l] = n
	}
	return res
}

var toClass = map[Name]*Class {
	FU: &Class {
		name: FU,
		leap: []V{newV(0, -1)},
		promote: TO,
	},
	TO: &Class {
		name: TO,
		leap: kinLeap,
		unpromote: FU,
	},
	KY: &Class {
		name: KY,
		ride: []V{newV(0, -1)},
		promote: NKY,
	},
	NKY: &Class {
		name: NKY,
		leap: kinLeap,
		unpromote: KY,
	},
	KE: &Class {
		name: KE,
		leap: []V{newV(1, -2), newV(-1, -2)},
		promote: NKE,
	},
	NKE: &Class {
		name: NKE,
		leap: kinLeap,
		unpromote: KE,
	},
	GI: &Class {
		name: GI,
		leap: []V{
			newV(1, -1), newV(0, -1), newV(-1, -1),
			//
			newV(1, 1), /*         */ newV(-1, 1),
		},
		promote: NGI,
	},
	NGI: &Class {
		name: NGI,
		leap: kinLeap,
		unpromote: GI,
	},
	KI: &Class {
		name: KI,
		leap: kinLeap,
	},
	KA: &Class {
		name: KA,
		ride: []V{
			newV(1, -1), newV(-1, -1),
			newV(1, 1), newV(-1, 1),
		},
		promote: UM,
	},
	UM: &Class {
		name: UM,
		leap: []V{newV(0, -1), newV(1, 0), newV(-1, 0), newV(0, 1)},
		ride: []V{
			newV(1, -1), newV(-1, -1),
			newV(1, 1), newV(-1, 1),
		},
		unpromote: KA,
	},
	HI: &Class {
		name: HI,
		ride: []V{newV(0, -1), newV(1, 0), newV(-1, 0), newV(0, 1)},
		promote: RY,
	},
	RY: &Class {
		name: RY,
		ride: []V{newV(0, -1), newV(1, 0), newV(-1, 0), newV(0, 1)},
		leap: []V{
			newV(1, -1), newV(-1, -1),
			newV(1, 1), newV(-1, 1),
		},
		unpromote: HI,
	},
	OU: &Class {
		name: OU,
		leap: []V{
			newV(1, -1), newV(0, -1), newV(-1, -1),
			newV(1, 0), /*         */ newV(-1, 0),
			newV(1, 1), newV(0, 1), newV(-1, 1),
		},
		king: true,
	},
}

func (n Name) class() *Class {
	return toClass[n]
}

func (n Name) letter() string {
	return toLetter[n]
}
