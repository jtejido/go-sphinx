package props

import (
	"os"

	"github.com/jtejido/go-sphinx/util"
)

/** Loads configuration from a TOML file */
type XMLLoader struct {
	file              *os.File
	rpdMap            map[string]*RawPropertyData
	globalProperties  map[string]string
	replaceDuplicates bool
}

/**
 * Creates a loader that will load from the given location
 *
 * @param url              the location to load
 * @param globalProperties the map of global properties
 * @param initRPD init raw property data
 * @param replaceDuplicates replace duplicates
 */
func NewXMLLoaderWithRPD(file *os.File, globalProperties map[string]string, initRPD map[string]*RawPropertyData, replaceDuplicates bool) *XMLLoader {
	ans := new(XMLLoader)
	ans.file = file
	ans.globalProperties = globalProperties
	ans.replaceDuplicates = replaceDuplicates
	if initRPD == nil {
		ans.rpdMap = make(map[string]*RawPropertyData)
	} else {
		ans.rpdMap = initRPD
	}
	return ans
}

/**
 * Creates a loader that will load from the given location
 *
 * @param url the location to load
 * @param globalProperties the map of global properties
 */
func NewXMLLoader(file *os.File, globalProperties map[string]string) *XMLLoader {
	return NewXMLLoaderWithRPD(file, globalProperties, nil, false)
}

/**
 * Loads a set of configuration data from the location
 *
 * @return a map keyed by component name containing RawPropertyData objects
 * @throws IOException if an I/O or parse error occurs
 */
func (l *XMLLoader) Load() (map[string]*RawPropertyData, error) {
	handler := NewConfigHandlerFromFile(l.rpdMap, l.globalProperties, l.replaceDuplicates, l.file)
	xr := util.NewParser(l.file, handler)
	if err := xr.Parse(); err != nil {
		return nil, err
	}

	return l.rpdMap, nil
}
