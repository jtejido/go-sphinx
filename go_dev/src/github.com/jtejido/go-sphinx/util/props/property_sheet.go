package props

const (
	COMP_LOG_LEVEL = "logLevel"
)

type PropertySheet struct {
	registeredProperties map[string]*S4PropWrapper
	propValues           map[string]any
	rawProps             map[string]any
	cm                   *ConfigurationManager
	owner                Configurable
	ownerClass           Configurable
	instanceName         string
}

/** @return the names of registered properties of this PropertySheet object. */
func (ps *PropertySheet) RegisteredProperties() []string {
	s := make([]string, len(ps.registeredProperties))
	var i int
	for k, _ := range ps.registeredProperties {
		s[i] = k
		i++
	}
	return s
}
