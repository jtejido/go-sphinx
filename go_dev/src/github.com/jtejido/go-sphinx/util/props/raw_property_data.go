package props

type RawPropertyData struct {
	name, className string
	properties      map[string]any
}

func NewRawPropertyData(name, className string) *RawPropertyData {
	return NewRawPropertyDataFromProps(name, className, make(map[string]any))
}

/**
 * Creates a raw property data item, using a given property map.
 *
 * @param name the name of the item
 * @param className the class name of the item
 * @param properties existing property map to use
 */
func NewRawPropertyDataFromProps(name, className string, properties map[string]any) *RawPropertyData {
	return &RawPropertyData{
		name:       name,
		className:  className,
		properties: properties,
	}
}

/** @return the className. */
func (d *RawPropertyData) ClassName() string {
	return d.className
}

/** @return the name. */
func (d *RawPropertyData) Name() string {
	return d.name
}

/** @return the properties. */
func (d *RawPropertyData) Properties() map[string]any {
	return d.properties
}

/**
 * Adds a new property with a {@code List<String>} value.
 *
 * @param propName the name of the property
 * @param propValue the value of the property
 */
func (d *RawPropertyData) Add(propName string, propValue string) {
	d.properties[propName] = propValue
}

/**
 * Adds a new property with a {@code List<String>} value.
 *
 * @param propName the name of the property
 * @param propValue the value of the property
 */
func (d *RawPropertyData) AddValues(propName string, propValue []string) {
	d.properties[propName] = propValue
}

/**
 * Determines if the map already contains an entry for a property.
 *
 * @param propName the property of interest
 * @return true if the map already contains this property
 */
func (d *RawPropertyData) Contains(propName string) (ok bool) {
	_, ok = d.properties[propName]
	return
}
