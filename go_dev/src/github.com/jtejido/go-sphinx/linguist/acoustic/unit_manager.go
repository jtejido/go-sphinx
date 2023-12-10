package acoustic

import "github.com/jtejido/go-sphinx/util"

const (
	SILENCE_NAME = "SIL"
	SILENCE_ID   = 1
)

var (
	SILENCE = NewUnit(SILENCE_NAME, true, SILENCE_ID)
)

type UnitManager struct {
	ciMap  map[string]*Unit
	nextID int
	logger util.Logger
}

func NewUnitManager(logger util.Logger) *UnitManager {
	return &UnitManager{
		ciMap: map[string]*Unit{
			SILENCE_NAME: SILENCE,
		},
		nextID: SILENCE_ID + 1,
		logger: logger,
	}
}

/**
 * Gets or creates a unit from the unit pool
 *
 * @param name    the name of the unit
 * @param filler  <code>true</code> if the unit is a filler unit
 * @param context the context for this unit
 * @return the unit
 */
func (um *UnitManager) UnitFromContext(name string, filler bool, context *Context) *Unit {
	unit := um.ciMap[name]
	if context == EMPTY_CONTEXT {
		if unit == nil {
			unit = NewUnit(name, filler, um.nextID)
			um.nextID++
			um.ciMap[name] = unit
			if um.logger != nil {
				um.logger.Infof("CI Unit: %s", unit.String())
			}
		}
	} else {
		unit = NewUnitFromContext(unit, filler, context)
	}
	return unit
}

/**
 * Gets or creates a unit from the unit pool
 *
 * @param name   the name of the unit
 * @param filler <code>true</code> if the unit is a filler unit
 * @return the unit
 */
func (um *UnitManager) Unit(name string, filler bool) *Unit {
	return um.UnitFromContext(name, filler, EMPTY_CONTEXT)
}

/**
 * Gets or creates a unit from the unit pool
 *
 * @param name the name of the unit
 * @return the unit
 */
func (um *UnitManager) UnitFromName(name string) *Unit {
	return um.UnitFromContext(name, false, EMPTY_CONTEXT)
}
