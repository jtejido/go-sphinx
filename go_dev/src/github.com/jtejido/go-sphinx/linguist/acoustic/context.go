package acoustic

var (
	EMPTY_CONTEXT = new(Context)
)

/** Represents  the context for a unit */
type Context struct {
}

/**
 * Checks to see if there is a partial match with the given context. For a simple context such as this we always
 * match.
 *
 * @param context the context to check
 * @return true if there is a partial match
 */
func (c *Context) IsPartialMatch(context *Context) bool {
	return true
}

/** Provides a string representation of a context */
func (c *Context) String() string {
	return ""
}

/**
 * Determines if an object is equal to this context
 *
 * @param o the object to check
 * @return true if the objects are equal
 */
func (c *Context) Equals(o any) bool {
	if c == o {
		return true
	} else if v, ok := o.(*Context); ok {
		return c.String() == v.String()
	} else {
		return false
	}
}
