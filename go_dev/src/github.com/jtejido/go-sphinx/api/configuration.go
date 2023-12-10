package api

// Represents common configuration options.
// This configuration is used by high-level recognition classes.
type Configuration struct {
	// path to acoustic model
	AcousticModelPath string
	// path to dictionary
	DictionaryPath string
	// path to the language model
	LanguageModelPath string
	// grammar path
	GrammarPath string
	// grammar name
	GrammarName string
	// The configured sample rate.
	SampleRate int
	// Whether fixed grammar should be used instead of language model.
	UseGrammar bool
}

func NewConfiguration() *Configuration {
	return &Configuration{
		SampleRate: 16000,
		UseGrammar: false,
	}
}
