package dictionary

const (
	DEFAULT_G2P_MODEL_PATH         string = ""
	DEFAULT_G2P_MAX_PRONUNCIATIONS        = 1

	// Spelling of the sentence start word.
	SENTENCE_START_SPELLING string = "<s>"

	// Spelling of the sentence end word.
	SENTENCE_END_SPELLING string = "</s>"

	// Spelling of the 'word' that marks a silence
	SILENCE_SPELLING string = "<sil>"
)

// Provides a generic interface to a dictionary. The dictionary is responsible for determining how a word is
// pronounced.
type Dictionary interface {
	// Returns a Word object based on the spelling and its classification. The behavior of this method is also affected
	// by the properties wordReplacement and g2pModel
	GetWord(text string) *Word

	// Returns the sentence start word.
	GetSentenceStartWord() *Word

	// Returns the sentence end word.
	GetSentenceEndWord() *Word

	// Returns the silence word.
	GetSilenceWord() *Word

	// Gets the set of all filler words in the dictionary
	GetFillerWords() []*Word

	// Allocates the dictionary
	Allocate()

	// Deallocates the dictionary
	Deallocate()
}
