package recognizer

import (
	"fmt"
	"io"
	"sync"

	"github.com/jtejido/go-sphinx/decoder"
	"github.com/jtejido/go-sphinx/instrumentation"
	"github.com/jtejido/go-sphinx/result"
)

// Called when the status has changed.
type StateListener func(State)

type block struct {
	try     func()
	catch   func(any)
	finally func()
}

func Throw(err any) {
	panic(err)
}

func (tcf block) do() {
	if tcf.finally != nil {
		defer tcf.finally()
	}
	if tcf.catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.catch(r)
			}
		}()
	}
	tcf.try()
}

type State int

const (
	// Defines the possible states of the recognizer
	DEALLOCATED State = iota
	ALLOCATING
	ALLOCATED
	READY
	RECOGNIZING
	DEALLOCATING
	ERROR
)

// The Sphinx-4 recognizer. This is the main entry point for Sphinx-4.
// Note that some Recognizer methods may throw panic if the recognizer is not in the proper state
type Recognizer struct {
	sync.Mutex
	name              string
	decoder           *decoder.Decoder
	currentState      State
	StateListenerFunc StateListener
	monitors          []instrumentation.Monitor
}

func NewDefaultRecognizer() *Recognizer {
	rec := new(Recognizer)
	rec.decoder = decoder.NewDefaultDecoder()
	rec.monitors = nil
	return rec
}

func NewRecognizer(decoder *decoder.Decoder, monitors []instrumentation.Monitor) *Recognizer {
	rec := new(Recognizer)
	rec.decoder = decoder
	rec.monitors = monitors

	return rec
}

// Performs recognition for the given number of input frames, or until a 'final' result is generated. This method
// should only be called when the recognizer is in the allocated state.
func (dr *Recognizer) Recognize(referenceText io.Reader) (result *result.Result) {

	dr.checkState(READY)

	block{
		try: func() {
			dr.setState(RECOGNIZING)
			result = dr.decoder.Decode(referenceText)
		},
		finally: func() {
			dr.setState(READY)
		},
	}.do()

	return result
}

// Checks to ensure that the recognizer is in the given state.
func (dr *Recognizer) checkState(desiredState State) {
	if dr.currentState != desiredState {
		panic(fmt.Sprintf("Expected state %d actual state %d", desiredState, dr.currentState))
	}
}

// sets the current state
func (dr *Recognizer) setState(newState State) {
	dr.currentState = newState
	if dr.StateListenerFunc != nil {
		dr.Lock()
		defer dr.Unlock()
		dr.StateListenerFunc(dr.currentState)
	}
}

// Allocate the resources needed for the recognizer. Note this method make take some time to complete. This method
// should only be called when the recognizer is in the deallocated state.
func (dr *Recognizer) Allocate() {
	dr.checkState(DEALLOCATED)
	dr.setState(ALLOCATING)
	dr.decoder.Allocate()
	dr.setState(ALLOCATED)
	dr.setState(READY)
}

// Deallocates the recognizer. This method should only be called if the recognizer is in the allocated state.
func (dr *Recognizer) Deallocate() {
	dr.checkState(READY)
	dr.setState(DEALLOCATING)
	dr.decoder.Deallocate()
	dr.setState(DEALLOCATED)
}

// Retrieves the recognizer state. This method can be called in any state.
func (dr *Recognizer) State() State {
	return dr.currentState
}

// Resets the monitors monitoring this recognizer
func (dr *Recognizer) ResetMonitors() {
	for _, listener := range dr.monitors {
		l, ok := listener.(instrumentation.Resetable)
		if ok {
			l.Reset()
		}
	}
}

func (dr *Recognizer) String() string {
	return fmt.Sprintf("Recognizer: %s  State: %d", dr.name, dr.currentState)
}
