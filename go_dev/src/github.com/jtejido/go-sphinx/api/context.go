package api

import (
	"io"
	"strconv"

	"github.com/jtejido/go-sphinx/frontend/util"
	"github.com/jtejido/go-sphinx/linguist/acoustic/tiedstate"
	"github.com/jtejido/go-sphinx/util/props"
)

type Context struct {
	configurationManager *props.ConfigurationManager
}

// Constructs builder that uses default XML configuration.
func NewDefaultContext(config *Configuration) *Context {
	return NewContext("default.config.yml", config)
}

// Constructs builder using user-supplied toml configuration.
func NewContext(path string, config *Configuration) *Context {
	ctx := new(Context)
	ctx.configurationManager = props.NewConfigurationManager(path)

	ctx.SetAcousticModel(config.AcousticModelPath)
	ctx.SetDictionary(config.DictionaryPath)

	if len(config.GrammarPath) > 0 && config.UseGrammar {
		ctx.SetGrammar(config.GrammarPath, config.GrammarName)
	}

	if len(config.LanguageModelPath) > 0 && !config.UseGrammar {
		ctx.SetLanguageModel(config.LanguageModelPath)
	}

	ctx.SetSampleRate(config.SampleRate)
	return ctx
}

// Sets acoustic model location.
//
// It also reads feat.params which should be located at the root of
// acoustic model and sets corresponding parameters of MelFrequencyFilterBank2 instance.
//
// Accepts path to directory with acoustic model files.
func (ctx *Context) SetAcousticModel(path string) {
	ctx.SetLocalProperty("acousticModelLoader->location", path)
	ctx.SetLocalProperty("dictionary->fillerPath", util.PathJoin(path, "noisedict"))
}

// Sets dictionary.
// Accepts path to directory with dictionary files.
func (ctx *Context) SetDictionary(path string) {
	ctx.SetLocalProperty("dictionary->dictionaryPath", path)
}

// Sets sampleRate.
// Accepts sample rate of the input stream.
func (ctx *Context) SetSampleRate(sampleRate int) {
	ctx.SetLocalProperty("dataSource->sampleRate", strconv.Itoa(sampleRate))
}

// Sets path to the grammar files.
//
// Enables static grammar and disables probabilistic language model.
// JSGF and GrXML formats are supported.
//
// Accepts path to the grammar files and name of the main grammar to use.
func (ctx *Context) SetGrammar(path, name string) {

	ctx.SetLocalProperty("jsgfGrammar->grammarLocation", path)
	ctx.SetLocalProperty("jsgfGrammar->grammarName", name)
	ctx.SetLocalProperty("flatLinguist->grammar", "jsgfGrammar")
	ctx.SetLocalProperty("decoder->searchManager", "simpleSearchManager")
}

// Sets path to the language model.
//
// Enables probabilistic language model and disables static grammar.
// Currently it supports ".lm", ".dmp" and ".bin" file formats.
//
// Accepts path to the language model file.
func (ctx *Context) SetLanguageModel(path string) {
	// if (path.endsWith(".lm")) {
	//     setLocalProperty("simpleNGramModel->location", path);
	//     setLocalProperty(
	//         "lexTreeLinguist->languageModel", "simpleNGramModel");
	// } else if (path.endsWith(".dmp")) {
	//     setLocalProperty("largeTrigramModel->location", path);
	//     setLocalProperty(
	//         "lexTreeLinguist->languageModel", "largeTrigramModel");
	// } else if (path.endsWith(".bin")) {
	//     setLocalProperty("trieNgramModel->location", path);
	//     setLocalProperty(
	//         "lexTreeLinguist->languageModel", "trieNgramModel");
	// } else {
	//     throw new IllegalArgumentException(
	//         "Unknown format extension: " + path);
	// }

	ctx.SetLocalProperty("trieNgramModel->location", path)
	ctx.SetLocalProperty("lexTreeLinguist->languageModel", "trieNgramModel")
}

// Sets byte stream as the speech source.
func (ctx *Context) SetSpeechSource(stream io.Reader, timeFrame *util.TimeFrame) {
	ds := ctx.GetInstance("dataSource").(util.StreamDataSource)
	ds.SetInputStream(stream, timeFrame)
	ctx.SetLocalProperty("trivialScorer->frontend", "liveFrontEnd")
}

// Sets property within a "component" tag in configuration.
//
// Use this method to alter "value" property of a "property" tag inside a
// "component" tag of the XML configuration.
// Accepts param name and value.
func (ctx *Context) SetLocalProperty(name, value string) {
	ctx.SetProperty(ctx.configurationManager, name, value)
}

// Sets property of a top-level "property" tag.
//
// Use this method to alter "value" property of a "property" tag whose
// parent is the root tag "config" of the XML configuration.
func (ctx *Context) SetGlobalProperty(name, value string) {
	ctx.configurationManager.SetGlobalProperty(name, value)
}

// Returns instance of the XML configuration by its class.
// note: cast it!
func (ctx *Context) GetInstance(c string) interface{} {
	return ctx.configurationManager.Lookup(c)
}

// Returns the Loader object used for loading the acoustic model.
func (ctx *Context) GetLoader() tiedstate.Loader {
	return ctx.configurationManager.Lookup("acousticModelLoader").(tiedstate.Loader)
}
