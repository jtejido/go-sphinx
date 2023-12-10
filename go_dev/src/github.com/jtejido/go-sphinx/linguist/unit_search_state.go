package linguist

import (
	"github.com/jtejido/go-sphinx/linguist/acoustic"
)

// Represents a unit state in a search space
type UnitSearchState interface {
	SearchState

	// Gets the unit
	GetUnit() acoustic.Unit
}
