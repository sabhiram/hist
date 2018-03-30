package types

type LineDesc struct {
	line    string // raw line / command
	comment string // editor comment / annotation
}

func NewLineDesc(l, c string) *LineDesc {
	return &LineDesc{
		line:    l,
		comment: c,
	}
}

func (ld *LineDesc) Line() string {
	if ld == nil {
		return ""
	}
	return ld.line
}

func (ld *LineDesc) Comment() string {
	if ld == nil {
		return ""
	}
	return ld.comment
}
