package instrumentation

/** Defines the interface for an object that is resetable */
type Resetable interface {
	/** Resets this component. Typically this is for components that keep track of statistics */
	Reset()
}
