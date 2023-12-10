package lextree

import (
	"github.com/jtejido/go-sphinx/linguist"
	"github.com/jtejido/go-sphinx/linguist/acoustic"
	"github.com/jtejido/go-sphinx/linguist/dictionary"
	"github.com/jtejido/go-sphinx/linguist/language/ngram"
	"github.com/jtejido/go-sphinx/linguist/tiedstate"
	"github.com/jtejido/go-sphinx/linguist/util"
	"github.com/jtejido/go-sphinx/utils"
)

const (
	DEFAULT_FULL_WORD_HISTORIES  bool    = true
	DEFAULT_CACHE_SIZE           int     = 0
	DEFAULT_ADD_FILLER_WORDS     bool    = false
	DEFAULT_GENERATE_UNIT_STATES bool    = false
	DEFAULT_WANT_UNIGRAM_SMEAR   bool    = true
	DEFAULT_UNIGRAM_SMEAR_WEIGHT float64 = 1.0
)

type LexTreeLinguist struct {

	// ----------------------------------
	// Subcomponents that are configured
	// by the property sheet
	// -----------------------------------
	languageModel ngram.LanguageModel
	acousticModel acoustic.AcousticModel
	dictionary    dictionary.Dictionary
	unitManager   acoustic.UnitManager

	// ------------------------------------
	// Data that is configured by the
	// property sheet
	// ------------------------------------
	addFillerWords, generateUnitStates, wantUnigramSmear       bool
	unigramSmearWeight                                         float64
	cacheEnabled                                               bool
	maxArcCacheSize                                            int
	languageWeight, logWordInsertionProbability                float64
	logUnitInsertionProbability, logFillerInsertionProbability float64
	logSilenceInsertionProbability, logOne                     float64

	// ------------------------------------
	// Data used for building and maintaining
	// the search graph
	// -------------------------------------
	sentenceEndWord        *dictionary.Word
	sentenceStartWordArray []*dictionary.Word
	searchGraph            linguist.SearchGraph
	hmmPool                *acoustic.HMMPool
	arcCache               *util.LRUCache
	maxDepth               int
	hmmTree                HMMTree
	cacheTrys, cacheHits   int
}

func NewLexTreeLinguist(acousticModel acoustic.AcousticModel, unitManager acoustic.UnitManager,
	languageModel language.LanguageModel, dictionary dictionary.Dictionary, fullWordHistories, wantUnigramSmear bool,
	wordInsertionProbability, silenceInsertionProbability, fillerInsertionProbability,
	unitInsertionProbability, languageWeight float64, addFillerWords, generateUnitStates bool,
	unigramSmearWeight float64, maxArcCacheSize int) *LexTreeLinguist {

	ltl := new(LexTreeLinguist)
	ltl.acousticModel = acousticModel
	ltl.unitManager = unitManager
	ltl.languageModel = languageModel
	ltl.dictionary = dictionary

	ltl.wantUnigramSmear = wantUnigramSmear
	ltl.logWordInsertionProbability = logMath.LinearToLog(wordInsertionProbability)
	ltl.logSilenceInsertionProbability = logMath.LinearToLog(silenceInsertionProbability)
	ltl.logFillerInsertionProbability = logMath.LinearToLog(fillerInsertionProbability)
	ltl.logUnitInsertionProbability = logMath.LinearToLog(unitInsertionProbability)
	ltl.languageWeight = languageWeight
	ltl.addFillerWords = addFillerWords
	ltl.generateUnitStates = generateUnitStates
	ltl.unigramSmearWeight = unigramSmearWeight
	ltl.maxArcCacheSize = maxArcCacheSize

	cacheEnabled := ltl.maxArcCacheSize > 0

	if cacheEnabled {
		ltl.arcCache = util.NewLRUCache(maxArcCacheSize)
	}

	return ltl
}

func NewDefaultLexTreeLinguist() *LexTreeLinguist {

	ltl := new(LexTreeLinguist)
	ltl.acousticModel = tiedstate.NewDefaultTiedStateAcousticModel()
	ltl.unitManager = acoustic.NewDefaultUnitManager()
	ltl.languageModel = ngram.NewDefaultSimpleNGramModel()
	ltl.dictionary = dictionary.NewDefaultTextDictionary()

	ltl.wantUnigramSmear = true
	ltl.logWordInsertionProbability = logMath.LinearToLog(0.1)
	ltl.logSilenceInsertionProbability = logMath.LinearToLog(0.1)
	ltl.logFillerInsertionProbability = logMath.LinearToLog(1e-2)
	ltl.logUnitInsertionProbability = logMath.LinearToLog(linguist.DEFAULT_UNIT_INSERTION_PROBABILITY)
	ltl.languageWeight = 8.0
	ltl.addFillerWords = true
	ltl.generateUnitStates = false
	ltl.unigramSmearWeight = 1.
	ltl.maxArcCacheSize = DEFAULT_CACHE_SIZE

	cacheEnabled := ltl.maxArcCacheSize > 0

	if cacheEnabled {
		ltl.arcCache = util.NewLRUCache(maxArcCacheSize)
	}

	return ltl
}

func (ltl *LexTreeLinguist) Allocate() {
	ltl.dictionary.Allocate()
	ltl.acousticModel.Allocate()
	ltl.languageModel.Allocate()
	ltl.compileGrammar()
}

func (ltl *LexTreeLinguist) Deallocate() {
	if ltl.acousticModel != nil {
		ltl.acousticModel.Deallocate()
	}

	if ltl.dictionary != nil {
		ltl.dictionary.Deallocate()
	}

	if ltl.languageModel != nil {
		ltl.languageModel.Deallocate()
	}

	ltl.hmmTree = nil
}

func (ltl LexTreeLinguist) GetSearchGraph() {
	return ltl.searchGraph
}

func (ltl *LexTreeLinguist) StartRecognition() {}

func (ltl *LexTreeLinguist) StopRecognition() {
	ltl.languageModel.OnUtteranceEnd()
}

func (ltl LexTreeLinguist) GetLanguageModel() ngram.LanguageModel {
	return ltl.languageModel
}

func (ltl LexTreeLinguist) GetDictionary() dictionary.Dictionary {
	return ltl.dictionary
}

// retrieves the initial language state
func (ltl LexTreeLinguist) getInitialSearchState() linguist.SearchState {
	node := ltl.hmmTree.GetInitialNode()

	if node == nil {
		panic("Language model has no entry for initial word <s>")
	}

	return NewLexTreeWordState(node, node.GetParent(), (NewWordSequenceByWordSlice(ltl.sentenceStartWordArray)).Trim(ltl.maxDepth-1), 0, ltl.logOne, ltl.logOne)
}

func (ltl *LexTreeLinguist) compileGrammar() {

	timer := utils.NewTimer("Compile")
	timer.Start()

	ltl.sentenceEndWord = dictionary.getSentenceEndWord()
	ltl.sentenceStartWordArray = []*dictionary.Word{ltl.dictionary.GetSentenceStartWord()}
	ltl.maxDepth = ltl.languageModel.GetMaxDepth()

	ltl.generateHmmTree()

	timer.Stop()
	// Now that we are all done, dump out some interesting
	// information about the process

	ltl.searchGraph = NewLexTreeSearchGraph(ltl.getInitialSearchState())
}

func (ltl *LexTreeLinguist) generateHmmTree() {
	ltl.hmmPool = acoustic.NewHMMPool(ltl.acousticModel, ltl.unitManager)
	ltl.hmmTree = NewHMMTree(ltl.hmmPool, ltl.dictionary, ltl.languageModel, ltl.addFillerWords, ltl.languageWeight)

	ltl.hmmPool.DumpInfo()
}
