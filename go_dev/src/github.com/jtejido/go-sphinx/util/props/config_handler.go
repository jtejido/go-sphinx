package props

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jtejido/go-sphinx/util"
)

type ConfigHandler struct {
	util.BaseHandler
	rpd          *RawPropertyData
	itemList     []string
	itemListName string
	curItem      string

	rpdMap           map[string]*RawPropertyData
	globalProperties map[string]string

	replaceDuplicates bool
	file              *os.File
}

func NewConfigHandlerFromFile(rpdMap map[string]*RawPropertyData, globalProperties map[string]string, replaceDuplicates bool, file *os.File) *ConfigHandler {
	return &ConfigHandler{
		rpdMap:            rpdMap,
		globalProperties:  globalProperties,
		replaceDuplicates: replaceDuplicates,
		file:              file,
	}
}

func NewConfigHandler(rpdMap map[string]*RawPropertyData, globalProperties map[string]string) *ConfigHandler {
	return NewConfigHandlerFromFile(rpdMap, globalProperties, false, nil)
}

func (h *ConfigHandler) StartElement(elem xml.StartElement) {
	if elem.Name.Local == "config" {
		// test if this configuration extends another one
		for _, attr := range elem.Attr {
			if attr.Name.Local == "extends" {
				extendedConfigName := attr.Value
				if extendedConfigName != "" {
					h.mergeConfigs(extendedConfigName, true)
					h.replaceDuplicates = true
				}
				break
			}
		}
	} else if elem.Name.Local == "include" {
		for _, attr := range elem.Attr {
			if attr.Name.Local == "file" {
				includeFileName := attr.Value
				h.mergeConfigs(includeFileName, false)
				break
			}
		}
	} else if elem.Name.Local == "extendwith" {
		for _, attr := range elem.Attr {
			if attr.Name.Local == "file" {
				includeFileName := attr.Value
				h.mergeConfigs(includeFileName, true)
				break
			}
		}
	} else if elem.Name.Local == "component" {
		var curComponent, curType string
		for _, attr := range elem.Attr {
			if attr.Name.Local == "name" {
				curComponent = attr.Value
			} else if attr.Name.Local == "type" {
				curType = attr.Value
			}
		}

		if h.rpdMap[curComponent] != nil && !h.replaceDuplicates {
			panic(fmt.Sprintf("duplicate definition for %s", curComponent))
		}
		h.rpd = NewRawPropertyData(curComponent, curType)
	} else if elem.Name.Local == "property" {
		var name, value string
		for _, attr := range elem.Attr {
			if attr.Name.Local == "name" {
				name = attr.Value
			} else if attr.Name.Local == "value" {
				value = attr.Value
			}
		}
		if len(elem.Attr) != 2 || name == "" || value == "" {
			panic("property element must only have 'name' and 'value' attributes")
		}
		if h.rpd == nil {
			// we are not in a component so add this to the global
			// set of symbols
			//                    String symbolName = "${" + name + "}"; // why should we warp the global props here
			h.globalProperties[name] = value
		} else if h.rpd.Contains(name) && !h.replaceDuplicates {
			panic(fmt.Sprintf("Duplicate property: %s", name))
		} else {
			h.rpd.Add(name, value)
		}
	} else if elem.Name.Local == "propertylist" {
		var itemListName string
		for _, attr := range elem.Attr {
			if attr.Name.Local == "name" {
				itemListName = attr.Value
				break
			}
		}

		if len(elem.Attr) != 1 || itemListName == "" {
			panic("list element must only have the 'name'  attribute")
		}
		h.itemList = make([]string, 0)
	} else if elem.Name.Local == "item" {
		if len(elem.Attr) != 0 {
			panic("unknown 'item' attribute")
		}
		h.curItem = ""
	} else {
		panic(fmt.Sprintf("Unknown element '%s'", elem.Name.Local))
	}
}

func (h *ConfigHandler) CharData(d xml.CharData) {
	if h.curItem != "" {
		h.curItem += string(d)
	}
}

func (h *ConfigHandler) EndElement(elem xml.EndElement) {
	if elem.Name.Local == "component" {
		h.rpdMap[h.rpd.Name()] = h.rpd
		h.rpd = nil
	} else if elem.Name.Local == "property" {
		// nothing to do
	} else if elem.Name.Local == "propertylist" {
		if h.rpd.Contains(h.itemListName) {
			panic(fmt.Sprintf("Duplicate property: %s", h.itemListName))
		} else {
			h.rpd.AddValues(h.itemListName, h.itemList)
			h.itemList = nil
		}
	} else if elem.Name.Local == "item" {
		h.itemList = append(h.itemList, strings.TrimSpace(h.curItem))
		h.curItem = ""
	}
}

func (h *ConfigHandler) mergeConfigs(configFileName string, replaceDuplicates bool) {
	basePath, err := filepath.Abs(h.file.Name())
	if err != nil {
		panic(fmt.Sprintf("Error getting absolute path: %v", err))
	}

	// Construct the file path for the configuration file
	configFilePath := filepath.Join(basePath, configFileName)

	configFile, err := os.Open(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("Error opening configuration file: %v", err))
	}
	defer configFile.Close()

	saxLoader := NewXMLLoaderWithRPD(configFile, h.globalProperties, h.rpdMap, replaceDuplicates)
	if _, err := saxLoader.Load(); err != nil {
		panic(fmt.Sprintf("Error loading file: %v", err))
	}

	logger := log.New(os.Stdout, "ConfigHandler: ", log.LstdFlags)
	logger.Printf("%s config: %s\n", getLoadAction(replaceDuplicates), configFile.Name())
}

func getParentDirectory(u *url.URL) string {
	path := u.Path
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}

func getFileURL(parent, fileName string) (*url.URL, error) {
	fileURL := fmt.Sprintf("%s/%s", parent, fileName)
	return url.Parse(fileURL)
}

func getLoadAction(replaceDuplicates bool) string {
	if replaceDuplicates {
		return "extending"
	}
	return "including"
}
