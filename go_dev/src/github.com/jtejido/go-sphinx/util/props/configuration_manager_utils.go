package props

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

const (
	/**
	 * A common property (used by all components) that sets the log level for the component.
	 */
	GLOBAL_COMMON_LOGLEVEL = "logLevel"

	/**
	 * The default file suffix of configuration files.
	 */
	CM_FILE_SUFFIX = ".sxl"
	FileScheme     = "file"
)

func GetURL(file *os.File) (*url.URL, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	path := fileInfo.Name()

	if !isWindowsDrivePath(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}

	if isWindowsDrivePath(path) {
		path = "/" + path
	}

	path = filepath.ToSlash(path)
	return &url.URL{
		Scheme: FileScheme,
		Path:   path,
	}, nil
}

// isWindowsDrivePath returns true if the file path is of the form used by Windows.
//
// We check if the path begins with a drive letter, followed by a ":".
func isWindowsDrivePath(path string) bool {
	if len(path) < 4 {
		return false
	}
	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}

// isWindowsDriveURI returns true if the file URI is of the format used by
// Windows URIs. The url.Parse package does not specially handle Windows paths
// (see https://golang.org/issue/6027). We check if the URI path has
// a drive prefix (e.g. "/C:"). If so, we trim the leading "/".
func isWindowsDriveURI(uri string) bool {
	if len(uri) < 4 {
		return false
	}
	return uri[0] == '/' && unicode.IsLetter(rune(uri[1])) && uri[2] == ':'
}

var jarPattern = regexp.MustCompile("(?i)resource:(.*)")

func ResourceToURL(location string) (*url.URL, error) {
	jarMatcher := jarPattern.FindStringSubmatch(location)
	if jarMatcher != nil {
		resourceName := jarMatcher[1]
		return url.Parse("file:" + resourceName)
	} else {
		if location[:5] != "file:" {
			location = "file:" + location
		}
		return url.Parse(location)
	}
}

// remark: the replacement of xml/sxl suffix is not necessary and just done to improve readability
func LogPrefix(cm *ConfigurationManager) string {
	if cm.ConfigFile() != nil {
		fileName := filepath.Base(cm.configFile.Name())
		return strings.TrimSuffix(strings.TrimSuffix(fileName, ".sxl"), ".xml") + "."

	}
	return "S4CM."
}

/**
 * Returns a map of all component-properties of this config-manager (including their associated property-sheets.
 *
 * @param cm configuration manager
 * @return map with properties
 */
func ListAllsPropNames(cm *ConfigurationManager) map[string][]*PropertySheet {
	allProps := make(map[string][]*PropertySheet)
	for _, configName := range cm.ComponentNames() {
		ps := cm.PropertySheet(configName)

		for _, propName := range ps.RegisteredProperties() {
			if _, ok := allProps[propName]; !ok {
				allProps[propName] = make([]*PropertySheet, 0)
			}
			allProps[propName] = append(allProps[propName], ps)
		}
	}

	return allProps
}

/**
 * Attempts to set the value of an arbitrary component-property. If the property-name is ambiguous  with respect to
 * the given <code>ConfiguratioManager</code> an extended syntax (componentName-&gt;propName) can be used to access the
 * property.
 * <p>
 * Beside component properties it is also possible to modify the class of a configurable, but this is only allowed if
 * the configurable under question has not been instantiated yet. Furthermore the user has to ensure to set all
 * mandatory component properties.
 * @param cm configuration manager
 * @param propName property to set
 * @param propValue value to set
 */
// func SetProperty( cm *ConfigurationManager,  propName,  propValue string) {
//         assert(propValue != "");
//        allProps := ListAllsPropNames(cm);
//         configurableNames := cm.ComponentNames();

//         if (!allProps.containsKey(propName) && !propName.contains("->") && !configurableNames.contains(propName))
//             throw new RuntimeException("No property or configurable '" + propName + "' in configuration '" + cm.getConfigURL() + "'!");

//         // if a configurable-class should be modified
//         if (configurableNames.contains(propName)) {
//             try {
//                 final Class<? extends Configurable> confClass = Class.forName(propValue).asSubclass(Configurable.class);
//                 ConfigurationManagerUtils.setClass(cm.getPropertySheet(propName), confClass);
//             } catch (ClassNotFoundException e) {
//                 throw new RuntimeException(e);
//             }

//             return;
//         }

//         if (!propName.contains("->") && allProps.get(propName).size() > 1) {
//             throw new RuntimeException("Property-name '" + propName + "' is ambiguous with respect to configuration '"
//                     + cm.getConfigURL() + "'. Use 'componentName->propName' to disambiguate your request.");
//         }

//         String componentName;

//         // if disambiguation syntax is used find the correct PS first
//         if (propName.contains("->")) {
//             String[] splitProp = propName.split("->");
//             componentName = splitProp[0];
//             propName = splitProp[1];
//         } else {
//             componentName = allProps.get(propName).get(0).getInstanceName();
//         }

//         setProperty(cm, componentName, propName, propValue);
//     }

/**
 * Configure the logger
 * @param cm Configuration manager
 */
func ConfigureLogger(cm *ConfigurationManager) {
	// apply the log level (if defined) for the root logger (because we're using package based logging now)
	cmPrefix := LogPrefix(cm)
	cmRootLogger := logrus.New()
	cmRootLogger.SetFormatter(&logrus.TextFormatter{
		FieldMap: logrus.FieldMap{
			"prefix": cmPrefix,
		},
	})

	//configureLogger(cmRootLogger)

	var level logrus.Level
	switch cm.GlobalProperty(GLOBAL_COMMON_LOGLEVEL) {
	case "info":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "warn":
		level = logrus.WarnLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		level = logrus.InfoLevel
	}

	cmRootLogger.SetLevel(level)
}

func assert(ok bool) {
	if !ok {
		panic("assert fail")
	}
}
