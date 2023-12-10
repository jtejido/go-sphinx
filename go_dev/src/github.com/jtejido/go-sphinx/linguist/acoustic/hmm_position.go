package acoustic

// HMMPosition type definition in Go
type HMMPosition int

const (
	BEGIN HMMPosition = iota
	END
	SINGLE
	INTERNAL
	UNDEFINED
)

var posByRep = map[rune]HMMPosition{
	'b': BEGIN,
	'e': END,
	's': SINGLE,
	'i': INTERNAL,
	'-': UNDEFINED,
}

// values function for HMMPosition in Go
func Values() []HMMPosition {
	return []HMMPosition{BEGIN, END, SINGLE, INTERNAL, UNDEFINED}
}

// String method for HMMPosition in Go
func (p HMMPosition) String() string {
	for rep, pos := range posByRep {
		if p == pos {
			return string(rep)
		}
	}
	return ""
}

// Lookup function for HMMPosition in Go
func Lookup(rep string) HMMPosition {
	if rep == "" {
		return UNDEFINED
	}
	return posByRep[rune(rep[0])]
}

// isWordEnd method for HMMPosition in Go
func (p HMMPosition) IsWordEnd() bool {
	return p == SINGLE || p == END
}

// isWordBeginning method for HMMPosition in Go
func (p HMMPosition) IsWordBeginning() bool {
	return p == SINGLE || p == BEGIN
}
