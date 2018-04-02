package types

type LineDesc struct {
	Line    string // raw line / command
	Comment string // editor comment / annotation
}

func NewLineDesc(l, c string) *LineDesc {
	return &LineDesc{
		Line:    l,
		Comment: c,
	}
}
