package props

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type ConfigurationManager struct {
	changeListeners  []ConfigurationChangeListener
	symbolTable      map[string]*PropertySheet
	rawPropertyMap   map[string]*RawPropertyData
	globalProperties map[string]string
	showCreations    bool
	configFile       *os.File
}

/**
 * Creates a new empty configuration manager. This constructor is only of use in cases when a system configuration
 * is created during runtime.
 */
func NewEmptyConfigurationManager() *ConfigurationManager {
	return new(ConfigurationManager)
}

/**
 * Creates a new configuration manager. Initial properties are loaded from the given URL. No need to keep the notion
 * of 'context' around anymore we will just pass around this property manager.
 *
 * @param configFileName The location of the configuration file.
 */
func NewConfigurationManager(configFileName string) (*ConfigurationManager, error) {
	file, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewConfigurationManagerFromFile(file)
}

/**
 * Creates a new configuration manager. Initial properties are loaded from the given URL. No need to keep the notion
 * of 'context' around anymore we will just pass around this property manager.
 *
 * @param url The location of the configuration file.
 */
func NewConfigurationManagerFromFile(file *os.File) (cm *ConfigurationManager, err error) {
	cm = new(ConfigurationManager)
	cm.configFile = file
	cm.globalProperties = make(map[string]string)
	cm.changeListeners = make([]ConfigurationChangeListener, 0)
	cm.symbolTable = make(map[string]*PropertySheet)

	cm.rawPropertyMap, err = NewXMLLoader(file, cm.globalProperties).Load()
	if err != nil {
		return nil, err
	}

	// ConfigurationManagerUtils.applySystemProperties(rawPropertyMap, globalProperties);
	ConfigureLogger(cm)

	// we can't configure the configuration manager with itself so we
	// do some of these configure items manually.
	showCreations := cm.globalProperties["showCreations"]
	if showCreations != "" {
		cm.showCreations = "true" == showCreations
	}

	return
}

/**
 * Returns the property sheet for the given object instance
 *
 * @param instanceName the instance name of the object
 * @return the property sheet for the object.
 */
func (cm *ConfigurationManager) PropertySheet(instanceName string) *PropertySheet {
	// if v,ok :=cm.symbolTable[instanceName];!ok {
	//     // if it is not in the symbol table, so construct
	//     // it based upon our raw property data
	//      rpd := cm.rawPropertyMap[instanceName]
	//     if (rpd != nil) {
	//          className := rpd.ClassName();
	//         try {
	//             Class<?> cls = Class.forName(className);

	//             // now load the property-sheet by using the class annotation
	//             PropertySheet propertySheet = new PropertySheet(cls.asSubclass(Configurable.class), instanceName, this, rpd);

	//             symbolTable.put(instanceName, propertySheet);

	//         } catch (ClassNotFoundException e) {
	//             System.err.println("class not found !" + e);
	//         } catch (ClassCastException e) {
	//             System.err.println("can not cast class !" + e);
	//         } catch (ExceptionInInitializerError e) {
	//             System.err.println("couldn't load class !" + e);
	//         }
	//     }
	// }

	return cm.symbolTable[instanceName]
}

func Lookup[V Configurable](cm *ConfigurationManager, instanceName string) (V, error) {
	// Apply all new properties to the model.
	instanceName = cm.StrippedComponentName(instanceName)
	ps = cm.PropertySheet(instanceName)

	if ps == nil {
		return nil, nil
	}

	if cm.showCreations {
		cm.RootLogger().config("Creating: " + instanceName)
	}

	return ps.Owner()
}

func (cm *ConfigurationManager) RootLogger() *logrus.Logger {
	return Logger.getLogger(ConfigurationManagerUtils.getLogPrefix(this))
	cmRootLogger := logrus.New()
	cmRootLogger.SetFormatter(&logrus.TextFormatter{
		FieldMap: logrus.FieldMap{
			"prefix": cmPrefix,
		},
	})
}

func (cm *ConfigurationManager) StrippedComponentName(propertyName string) string {
	assert(propertyName != "")

	for {
		if !strings.HasPrefix(propertyName, "$") {
			break
		}
		propertyName = cm.globalProperties[StripGlobalSymbol(propertyName)]
	}

	return propertyName
}

/**
 * Returns a global property.
 *
 * @param propertyName The name of the global property or <code>null</code> if no such property exists
 * @return a global property
 */
func (cm *ConfigurationManager) GlobalProperty(propertyName string) string {
	// propertyName = propertyName.startsWith("$") ? propertyName : "${" + propertyName + "}";
	globProp := cm.globalProperties[propertyName]
	if globProp != "" {
		return globProp
	}
	return ""
}

/**
 * Returns all names of configurables registered to this instance. The resulting set includes instantiated and
 * non-instantiated components.
 *
 * @return all component named registered to this instance of <code>ConfigurationManager</code>
 */
func (cm *ConfigurationManager) ComponentNames() []string {
	s := make([]string, len(cm.rawPropertyMap))
	var i int
	for k, _ := range cm.rawPropertyMap {
		s[i] = k
		i++
	}
	return s
}

func (cm *ConfigurationManager) ConfigFile() *os.File {
	return cm.configFile
}
