package util

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
)

/** A class that provides a mechanism for tokenizing a stream */
type ExtendedStreamTokenizer struct {
	path        string
	st          *StreamTokenizer
	reader      *bufio.Reader
	atEOF       bool
	putbackList []string
}

// NewExtendedStreamTokenizer creates and returns a new ExtendedStreamTokenizer
func NewExtendedStreamTokenizerFromReader(reader io.Reader, commentChar int, eolIsSignificant bool) (*ExtendedStreamTokenizer, error) {
	rr := bufio.NewReader(reader)
	st, err := NewStreamTokenizerFromReader(rr)
	if err != nil {
		return nil, err
	}
	st.ResetSyntax()
	st.WhitespaceChars(0, 32)
	st.WordChars(33, 255)
	st.EOLIsSignificant(eolIsSignificant)
	st.CommentChar(commentChar)
	return &ExtendedStreamTokenizer{
		reader:      rr,
		st:          st,
		putbackList: make([]string, 0),
	}, nil
}

/**
 * Specifies that all the characters between low and hi incluseive are whitespace characters
 *
 * @param low the low end of the range
 * @param hi  the high end of the range
 */
func (t *ExtendedStreamTokenizer) WhitespaceChars(low, hi int) {
	t.st.WhitespaceChars(low, hi)
}

/**
 * Specified that the character argument starts a single-line comment. All characters from the comment character to
 * the end of the line are ignored by this stream tokenizer.
 *
 * @param ch the comment character
 */
func (t *ExtendedStreamTokenizer) CommentChar(ch int) {
	t.st.CommentChar(ch)
}

/**
 * Gets the next word from the tokenizer
 *
 * @return the next word
 * @throws StreamCorruptedException if the word does not match
 * @throws IOException              if an error occurs while loading the data
 */
func (t *ExtendedStreamTokenizer) GetString() (string, error) {
	if len(t.putbackList) > 0 {
		i := len(t.putbackList) - 1
		s := t.putbackList[i]
		t.putbackList = append(t.putbackList[:i], t.putbackList[i+1:]...)
		return s, nil
	} else {
		_, err := t.st.NextToken()
		if err != nil {
			return "", err
		}
		if t.st.ttype == TT_EOF {
			t.atEOF = true
		}
		if t.st.ttype != TT_WORD &&
			t.st.ttype != TT_EOL &&
			t.st.ttype != TT_EOF {
			return "", t.corrupt("word expected but not found")
		}
		if t.st.ttype == TT_EOL ||
			t.st.ttype == TT_EOF {
			return "", nil
		} else {
			return t.st.sval, nil
		}
	}
}

/**
 * Puts a string back, the next get will return this string
 *
 * @param string the string to unget
 */
func (t *ExtendedStreamTokenizer) Unget(str string) {
	t.putbackList = append(t.putbackList, str)
}

/**
 * Determines if the stream is at the end of file
 *
 * @return true if the stream is at EOF
 */
func (t *ExtendedStreamTokenizer) IsEOF() bool {
	return t.atEOF
}

/**
 * Gets the current line number
 *
 * @return the line number
 */
func (t *ExtendedStreamTokenizer) LineNumber() int {
	return t.st.Lineno()
}

/**
 * Loads a word from the tokenizer and ensures that it matches 'expecting'
 *
 * @param expecting the word read must match this
 * @throws StreamCorruptedException if the word does not match
 * @throws IOException              if an error occurs while loading the data
 */
func (t *ExtendedStreamTokenizer) ExpectString(expecting string) error {
	line, err := t.GetString()
	if err != nil {
		return err
	}
	if line != expecting {
		return t.corrupt(fmt.Sprintf("error matching expected string '%s' in line: '%s'", expecting, line))
	}

	return nil
}

/**
 * Loads an integer  from the tokenizer and ensures that it matches 'expecting'
 *
 * @param name      the name of the value
 * @param expecting the word read must match this
 * @throws StreamCorruptedException if the word does not match
 * @throws IOException              if an error occurs while loading the data
 */
func (t *ExtendedStreamTokenizer) ExpectInt(name string, expecting int) error {
	val, err := t.GetInt(name)
	if err != nil {
		return err
	}
	if val != expecting {
		return t.corrupt(fmt.Sprintf("Expecting integer %d", expecting))
	}

	return nil
}

/**
 * gets an integer from the tokenizer stream
 *
 * @param name the name of the parameter (for error reporting)
 * @return the next word in the stream as an integer
 * @throws StreamCorruptedException if the next value is not a
 * @throws IOException              if an error occurs while loading the data number
 */
func (t *ExtendedStreamTokenizer) GetInt(name string) (iVal int, err error) {
	val, err := t.GetString()
	if err != nil {
		return iVal, t.corrupt(fmt.Sprintf("while parsing int %s", name))
	}
	iVal, err = strconv.Atoi(val)
	if err != nil {
		return iVal, t.corrupt(fmt.Sprintf("while parsing int %s", name))
	}

	return
}

/**
 * gets a double from the tokenizer stream
 *
 * @param name the name of the parameter (for error reporting)
 * @return the next word in the stream as a double
 * @throws StreamCorruptedException if the next value is not a
 * @throws IOException              if an error occurs while loading the data number
 */
func (t *ExtendedStreamTokenizer) GetFloat64(name string) (dVal float64, err error) {
	val, err := t.GetString()
	if err != nil {
		return dVal, t.corrupt(fmt.Sprintf("while parsing double %s", name))
	}

	if val == "inf" {
		return math.Inf(0), nil
	}

	dVal, err = strconv.ParseFloat(val, 64)
	if err != nil {
		return 0.0, t.corrupt(fmt.Sprintf("while parsing double %s", name))
	}

	return
}

/**
* gets a float from the tokenizer stream
*
* @param name the name of the parameter (for error reporting)
* @return the next word in the stream as a float
* @throws StreamCorruptedException if the next value is not a
* @throws IOException              if an error occurs while loading the data number
 */
func (t *ExtendedStreamTokenizer) GetFloat32(name string) (float32, error) {
	val, err := t.GetString()
	if err != nil {
		return 0.0, t.corrupt(fmt.Sprintf("while parsing float %s", name))
	}

	if val == "inf" {
		return float32(math.Inf(0)), nil
	}

	fVal, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0.0, t.corrupt(fmt.Sprintf("while parsing float %s", name))
	}

	return float32(fVal), nil

}

/**
 * gets a optional float from the tokenizer stream. If a float is not present, the default is returned
 *
 * @param name         the name of the parameter (for error reporting)
 * @param defaultValue the default value
 * @return the next word in the stream as a float
 * @throws StreamCorruptedException if the next value is not a
 * @throws IOException              if an error occurs while loading the data number
 */
func (t *ExtendedStreamTokenizer) GetFloat32WithDefault(name string, defaultValue float32) (float32, error) {
	val, err := t.GetString()
	if err != nil {
		return 0.0, t.corrupt(fmt.Sprintf("while parsing float %s", name))
	}

	if val == "" {
		return defaultValue, nil
	}

	if val == "inf" {
		return float32(math.Inf(0)), nil
	}

	fVal, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0.0, t.corrupt(fmt.Sprintf("while parsing float %s", name))
	}

	return float32(fVal), nil
}

/**
 * Skip any carriage returns.
 *
 * @throws IOException if an error occurs while reading data from the stream.
 */
func (t *ExtendedStreamTokenizer) SkipWhite() error {
	for !t.IsEOF() {
		token, err := t.GetString()
		if err != nil {
			return err
		}

		if token != "" {
			t.Unget(token)
			break
		}
	}

	return nil
}

func (t *ExtendedStreamTokenizer) corrupt(msg string) error {
	return &streamCorruptError{
		msg: fmt.Sprintf("%s at line %d in file %s", msg, t.st.Lineno(), t.path),
	}
}

type streamCorruptError struct {
	msg string
}

func (e *streamCorruptError) Error() string {
	return e.msg
}
