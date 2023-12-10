package acoustic

import "reflect"

var EMPTY_ARRAY = make([]*Unit, 0)

/** Represents a unit of speech. Units may represent phones, words or any other suitable unit */

type Unit struct {
	name            string
	filler, silence bool
	baseID          int
	baseUnit        *Unit
	context         *Context
	key             string
}

/**
 * Constructs a context independent unit. Constructors are package private, use the UnitManager to create and access
 * units.
 *
 * @param name   the name of the unit
 * @param filler <code>true</code> if the unit is a filler unit
 * @param id     the base id for the unit
 */
func NewUnit(name string, filler bool, id int) *Unit {
	this := new(Unit)
	this.name = name
	this.filler = filler
	this.silence = name == SILENCE_NAME
	this.baseID = id
	this.baseUnit = this
	this.context = EMPTY_CONTEXT
	return this
}

/**
 * Constructs a context dependent unit. Constructors are package private, use the UnitManager to create and access
 * units.
 *
 * @param baseUnit the base id for the unit
 * @param filler   <code>true</code> if the unit is a filler unit
 * @param context  the context for this unit
 */
func NewUnitFromContext(baseUnit *Unit, filler bool, context *Context) *Unit {
	this := new(Unit)
	this.name = baseUnit.Name()
	this.filler = filler
	this.silence = this.name == SILENCE_NAME
	this.baseID = baseUnit.BaseID()
	this.baseUnit = baseUnit
	this.context = context
	return this
}

/**
 * Gets the name for this unit
 *
 * @return the name for this unit
 */
func (u *Unit) Name() string {
	return u.name
}

/**
 * Determines if this unit is a filler unit
 *
 * @return <code>true</code> if the unit is a filler unit
 */
func (u *Unit) IsFiller() bool {
	return u.filler
}

/**
 * Determines if this unit is the silence unit
 *
 * @return true if the unit is the silence unit
 */
func (u *Unit) IsSilence() bool {
	return u.silence
}

/**
 * Gets the base ID for this unit
 *
 * @return the id
 */
func (u *Unit) BaseID() int {
	return u.baseID
}

/**
 * Gets the  base unit associated with this HMM
 *
 * @return the unit associated with this HMM
 */
func (u *Unit) BaseUnit() *Unit {
	return u.baseUnit
}

/**
 * Returns the context for this unit
 *
 * @return the context for this unit (or null if context independent)
 */
func (u *Unit) Context() *Context {
	return u.context
}

/**
 * Determines if this unit is context dependent
 *
 * @return true if the unit is context dependent
 */
func (u *Unit) IsContextDependent() bool {
	return u.Context() != EMPTY_CONTEXT
}

/** gets the key for this unit
 * @return the key
 */
func (u *Unit) Key() string {
	return u.String()
}

/**
 * Checks to see of an object is equal to this unit
 *
 * @param o the object to check
 * @return true if the objects are equal
 */
func (u *Unit) Equals(o any) bool {
	if u == o {
		return true
	} else if v, ok := o.(*Unit); ok {
		return u.Key() == v.Key()
	} else {
		return false
	}
}

/**
 * Converts to a string
 *
 * @return string version
 */
func (u *Unit) String() string {
	if u.key == "" {
		if u.context == EMPTY_CONTEXT {
			if u.filler {
				u.key = "*"
			} else {
				u.key = ""
			}
			u.key += u.name
		} else {

			if u.filler {
				u.key = "*"
			} else {
				u.key = ""
			}
			u.key += u.name + "[" + u.context.String() + "]"
		}
	}
	return u.key
}

/**
 * Checks to see if the given unit with associated contexts is a partial match for this unit.   Zero, One or both
 * contexts can be null. A null context matches any context
 *
 * @param name    the name of the unit
 * @param context the  context to match against
 * @return true if this unit matches the name and non-null context
 */
func (u *Unit) IsPartialMatch(name string, context *Context) bool {
	return u.Name() == name && context.IsPartialMatch(u.context)
}

/**
 * Creates and returns an empty context with the given size. The context is padded with SIL filler
 *
 * @param size the size of the context
 * @return the context
 */

func (u *Unit) EmptyContext(size int) []*Unit {
	context := make([]*Unit, size)
	for i := 0; i < size; i++ {
		context[i] = SILENCE
	}

	return context
}

/**
 * Checks to see that there is 100% overlap in the given contexts
 *
 * @param a context to check for a match
 * @param b context to check for a match
 * @return <code>true</code> if the contexts match
 */
func (u *Unit) IsContextMatch(a, b []*Unit) bool {
	if a == nil || b == nil {
		return reflect.DeepEqual(a, b)
	} else if len(a) != len(b) {
		return false
	} else {
		for i := 0; i < len(a); i++ {
			if a[i].Name() != b[i].Name() {
				return false
			}
		}
		return true
	}
}
