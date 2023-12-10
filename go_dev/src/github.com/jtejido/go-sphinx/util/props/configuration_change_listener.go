package props

/**
 * Describes all methods necessary to process change events of a <code>ConfigurationManager</code>.
 *
 * @see edu.cmu.sphinx.util.props.ConfigurationManager
 */
type ConfigurationChangeListener interface {

	/**
	 * Called if the configuration of a registered compoenent named <code>configurableName</code> was changed.
	 *
	 * @param configurableName The name of the changed configurable.
	 * @param propertyName     The name of the property which was changed
	 * @param cm               The <code>ConfigurationManager</code>-instance this component is registered to
	 */
	ConfigurationChanged(configurableName, propertyName string, cm *ConfigurationManager)

	/**
	 * Called if a new compoenent defined by <code>ps</code> was registered to the ConfigurationManager
	 * <code>cm</code>.
	 * @param cm               Configuration manager
	 * @param ps               Property sheet
	 */
	ComponentAdded(cm ConfigurationManager, ps *PropertySheet)

	/**
	 * Called if a compoenent defined by <code>ps</code> was unregistered (removed) from the ConfigurationManager
	 * <code>cm</code>.
	 * @param cm               Configuration manager
	 * @param ps               Property sheet
	 */
	ComponentRemoved(cm ConfigurationManager, ps *PropertySheet)

	/**
	 * Called if a compoenent was renamed.
	 * @param cm               Configuration manager
	 * @param ps               Property sheet
	 * @param oldName          Old name
	 */
	ComponentRenamed(cm ConfigurationManager, ps *PropertySheet, oldName string)
}
