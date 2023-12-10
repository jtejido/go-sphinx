package util

import (
	"fmt"
	"io"
	"math"
	"strings"
)

const (
	NEED_CHAR          = math.MaxInt
	SKIP_LF            = math.MaxInt - 1
	CT_WHITESPACE byte = 1
	CT_DIGIT      byte = 2
	CT_ALPHA      byte = 4
	CT_QUOTE      byte = 8
	CT_COMMENT    byte = 16
	/**
	 * A constant indicating that the end of the stream has been read.
	 */
	TT_EOF int = -1

	/**
	 * A constant indicating that the end of the line has been read.
	 */
	TT_EOL int = '\n'

	/**
	 * A constant indicating that a number token has been read.
	 */
	TT_NUMBER int = -2

	/**
	 * A constant indicating that a word token has been read.
	 */
	TT_WORD int = -3

	/* A constant indicating that no token has been read, used for
	 * initializing ttype.  FIXME This could be made public and
	 * made available as the part of the API in a future release.
	 */
	TT_NOTHING int = -4
)

type StreamTokenizer struct {
	/* Only one of these will be non-null */
	reader io.Reader
	buf    []rune

	/**
	 * The next character to be considered by the nextToken method.  May also
	 * be NEED_CHAR to indicate that a new character should be read, or SKIP_LF
	 * to indicate that a new character should be read and, if it is a '\n'
	 * character, it should be discarded and a second new character should be
	 * read.
	 */
	peekc               int
	pushedBack          bool
	forceLower          bool
	eolIsSignificantP   bool
	slashSlashCommentsP bool
	slashStarCommentsP  bool
	/** The line number of the last token read */
	lineno int
	ctype  []byte
	/**
	 * After a call to the <code>nextToken</code> method, this field
	 * contains the type of the token just read. For a single character
	 * token, its value is the single character, converted to an integer.
	 * For a quoted string token, its value is the quote character.
	 * Otherwise, its value is one of the following:
	 * <ul>
	 * <li><code>TT_WORD</code> indicates that the token is a word.
	 * <li><code>TT_NUMBER</code> indicates that the token is a number.
	 * <li><code>TT_EOL</code> indicates that the end of line has been read.
	 *     The field can only have this value if the
	 *     <code>eolIsSignificant</code> method has been called with the
	 *     argument <code>true</code>.
	 * <li><code>TT_EOF</code> indicates that the end of the input stream
	 *     has been reached.
	 * </ul>
	 * <p>
	 * The initial value of this field is -4.
	 *
	 * @see     java.io.StreamTokenizer#eolIsSignificant(boolean)
	 * @see     java.io.StreamTokenizer#nextToken()
	 * @see     java.io.StreamTokenizer#quoteChar(int)
	 * @see     java.io.StreamTokenizer#TT_EOF
	 * @see     java.io.StreamTokenizer#TT_EOL
	 * @see     java.io.StreamTokenizer#TT_NUMBER
	 * @see     java.io.StreamTokenizer#TT_WORD
	 */
	ttype int

	/**
	 * If the current token is a word token, this field contains a
	 * string giving the characters of the word token. When the current
	 * token is a quoted string token, this field contains the body of
	 * the string.
	 * <p>
	 * The current token is a word when the value of the
	 * <code>ttype</code> field is <code>TT_WORD</code>. The current token is
	 * a quoted string token when the value of the <code>ttype</code> field is
	 * a quote character.
	 * <p>
	 * The initial value of this field is null.
	 *
	 * @see     java.io.StreamTokenizer#quoteChar(int)
	 * @see     java.io.StreamTokenizer#TT_WORD
	 * @see     java.io.StreamTokenizer#ttype
	 */
	sval string

	/**
	 * If the current token is a number, this field contains the value
	 * of that number. The current token is a number when the value of
	 * the <code>ttype</code> field is <code>TT_NUMBER</code>.
	 * <p>
	 * The initial value of this field is 0.0.
	 *
	 * @see     java.io.StreamTokenizer#TT_NUMBER
	 * @see     java.io.StreamTokenizer#ttype
	 */
	nval float64
}

/** Private constructor that initializes everything except the streams. */
func NewDefaultStreamTokenizer() *StreamTokenizer {
	this := new(StreamTokenizer)
	this.peekc = NEED_CHAR
	this.buf = make([]rune, 20)
	this.ctype = make([]byte, 256)
	this.lineno = 1
	this.WordChars('a', 'z')
	this.WordChars('A', 'Z')
	this.WordChars(128+32, 255)
	this.WhitespaceChars(0, ' ')
	this.CommentChar('/')
	this.QuoteChar('"')
	this.QuoteChar('\'')
	this.ParseNumbers()
	return this
}

/**
 * Create a tokenizer that parses the given character stream.
 *
 * @param r  a Reader object providing the input stream.
 * @since   JDK1.1
 */
func NewStreamTokenizerFromReader(r io.Reader) (*StreamTokenizer, error) {
	this := NewDefaultStreamTokenizer()
	if r == nil {
		return nil, fmt.Errorf("reader cannot be nil")
	}
	this.reader = r
	return this, nil
}

/**
 * Resets this tokenizer's syntax table so that all characters are
 * "ordinary." See the <code>ordinaryChar</code> method
 * for more information on a character being ordinary.
 *
 * @see     java.io.StreamTokenizer#ordinaryChar(int)
 */
func (t *StreamTokenizer) ResetSyntax() {
	for i := range t.ctype {
		t.ctype[i] = 0
	}
}

/**
 * Specifies that all characters <i>c</i> in the range
 * <code>low&nbsp;&lt;=&nbsp;<i>c</i>&nbsp;&lt;=&nbsp;high</code>
 * are word constituents. A word token consists of a word constituent
 * followed by zero or more word constituents or number constituents.
 *
 * @param   low   the low end of the range.
 * @param   hi    the high end of the range.
 */
func (t *StreamTokenizer) WordChars(low, hi int) {
	if low < 0 {
		low = 0
	}
	if hi >= len(t.ctype) {
		hi = len(t.ctype) - 1
	}
	for low <= hi {
		t.ctype[low] |= CT_ALPHA
		low++
	}
}

/**
 * Specifies that all characters <i>c</i> in the range
 * <code>low&nbsp;&lt;=&nbsp;<i>c</i>&nbsp;&lt;=&nbsp;high</code>
 * are white space characters. White space characters serve only to
 * separate tokens in the input stream.
 *
 * <p>Any other attribute settings for the characters in the specified
 * range are cleared.
 *
 * @param   low   the low end of the range.
 * @param   hi    the high end of the range.
 */
func (t *StreamTokenizer) WhitespaceChars(low, hi int) {
	if low < 0 {
		low = 0
	}
	if hi >= len(t.ctype) {
		hi = len(t.ctype) - 1
	}
	for low <= hi {
		t.ctype[low] = CT_WHITESPACE
		low++
	}
}

/**
 * Specifies that all characters <i>c</i> in the range
 * <code>low&nbsp;&lt;=&nbsp;<i>c</i>&nbsp;&lt;=&nbsp;high</code>
 * are "ordinary" in this tokenizer. See the
 * <code>ordinaryChar</code> method for more information on a
 * character being ordinary.
 *
 * @param   low   the low end of the range.
 * @param   hi    the high end of the range.
 * @see     java.io.StreamTokenizer#ordinaryChar(int)
 */
func (t *StreamTokenizer) OrdinaryChars(low, hi int) {
	if low < 0 {
		low = 0
	}
	if hi >= len(t.ctype) {
		hi = len(t.ctype) - 1
	}
	for low <= hi {
		t.ctype[low] = 0
		low++
	}
}

/**
 * Specifies that the character argument is "ordinary"
 * in this tokenizer. It removes any special significance the
 * character has as a comment character, word component, string
 * delimiter, white space, or number character. When such a character
 * is encountered by the parser, the parser treats it as a
 * single-character token and sets <code>ttype</code> field to the
 * character value.
 *
 * <p>Making a line terminator character "ordinary" may interfere
 * with the ability of a <code>StreamTokenizer</code> to count
 * lines. The <code>lineno</code> method may no longer reflect
 * the presence of such terminator characters in its line count.
 *
 * @param   ch   the character.
 * @see     java.io.StreamTokenizer#ttype
 */
func (t *StreamTokenizer) OrdinaryChar(ch int) {
	if ch >= 0 && ch < len(t.ctype) {
		t.ctype[ch] = 0
	}
}

/**
 * Specified that the character argument starts a single-line
 * comment. All characters from the comment character to the end of
 * the line are ignored by this stream tokenizer.
 *
 * <p>Any other attribute settings for the specified character are cleared.
 *
 * @param   ch   the character.
 */
func (t *StreamTokenizer) CommentChar(ch int) {
	if ch >= 0 && ch < len(t.ctype) {
		t.ctype[ch] = CT_COMMENT
	}
}

/**
 * Specifies that matching pairs of this character delimit string
 * constants in this tokenizer.
 * <p>
 * When the <code>nextToken</code> method encounters a string
 * constant, the <code>ttype</code> field is set to the string
 * delimiter and the <code>sval</code> field is set to the body of
 * the string.
 * <p>
 * If a string quote character is encountered, then a string is
 * recognized, consisting of all characters after (but not including)
 * the string quote character, up to (but not including) the next
 * occurrence of that same string quote character, or a line
 * terminator, or end of file. The usual escape sequences such as
 * <code>"&#92;n"</code> and <code>"&#92;t"</code> are recognized and
 * converted to single characters as the string is parsed.
 *
 * <p>Any other attribute settings for the specified character are cleared.
 *
 * @param   ch   the character.
 * @see     java.io.StreamTokenizer#nextToken()
 * @see     java.io.StreamTokenizer#sval
 * @see     java.io.StreamTokenizer#ttype
 */
func (t *StreamTokenizer) QuoteChar(ch int) {
	if ch >= 0 && ch < len(t.ctype) {
		t.ctype[ch] = CT_QUOTE
	}
}

/**
 * Specifies that numbers should be parsed by this tokenizer. The
 * syntax table of this tokenizer is modified so that each of the twelve
 * characters:
 * <blockquote><pre>
 *      0 1 2 3 4 5 6 7 8 9 . -
 * </pre></blockquote>
 * <p>
 * has the "numeric" attribute.
 * <p>
 * When the parser encounters a word token that has the format of a
 * double precision floating-point number, it treats the token as a
 * number rather than a word, by setting the <code>ttype</code>
 * field to the value <code>TT_NUMBER</code> and putting the numeric
 * value of the token into the <code>nval</code> field.
 *
 * @see     java.io.StreamTokenizer#nval
 * @see     java.io.StreamTokenizer#TT_NUMBER
 * @see     java.io.StreamTokenizer#ttype
 */
func (t *StreamTokenizer) ParseNumbers() {
	for i := '0'; i <= '9'; i++ {
		t.ctype[i] |= CT_DIGIT
	}
	t.ctype['.'] |= CT_DIGIT
	t.ctype['-'] |= CT_DIGIT
}

/**
 * Determines whether or not ends of line are treated as tokens.
 * If the flag argument is true, this tokenizer treats end of lines
 * as tokens; the <code>nextToken</code> method returns
 * <code>TT_EOL</code> and also sets the <code>ttype</code> field to
 * this value when an end of line is read.
 * <p>
 * A line is a sequence of characters ending with either a
 * carriage-return character (<code>'&#92;r'</code>) or a newline
 * character (<code>'&#92;n'</code>). In addition, a carriage-return
 * character followed immediately by a newline character is treated
 * as a single end-of-line token.
 * <p>
 * If the <code>flag</code> is false, end-of-line characters are
 * treated as white space and serve only to separate tokens.
 *
 * @param   flag   <code>true</code> indicates that end-of-line characters
 *                 are separate tokens; <code>false</code> indicates that
 *                 end-of-line characters are white space.
 * @see     java.io.StreamTokenizer#nextToken()
 * @see     java.io.StreamTokenizer#ttype
 * @see     java.io.StreamTokenizer#TT_EOL
 */
func (t *StreamTokenizer) EOLIsSignificant(flag bool) {
	t.eolIsSignificantP = flag
}

/**
 * Determines whether or not the tokenizer recognizes C-style comments.
 * If the flag argument is <code>true</code>, this stream tokenizer
 * recognizes C-style comments. All text between successive
 * occurrences of <code>/*</code> and <code>*&#47;</code> are discarded.
 * <p>
 * If the flag argument is <code>false</code>, then C-style comments
 * are not treated specially.
 *
 * @param   flag   <code>true</code> indicates to recognize and ignore
 *                 C-style comments.
 */
func (t *StreamTokenizer) SlashStarComments(flag bool) {
	t.slashStarCommentsP = flag
}

/**
 * Determines whether or not the tokenizer recognizes C++-style comments.
 * If the flag argument is <code>true</code>, this stream tokenizer
 * recognizes C++-style comments. Any occurrence of two consecutive
 * slash characters (<code>'/'</code>) is treated as the beginning of
 * a comment that extends to the end of the line.
 * <p>
 * If the flag argument is <code>false</code>, then C++-style
 * comments are not treated specially.
 *
 * @param   flag   <code>true</code> indicates to recognize and ignore
 *                 C++-style comments.
 */
func (t *StreamTokenizer) SlashSlashComments(flag bool) {
	t.slashSlashCommentsP = flag
}

/**
 * Determines whether or not word token are automatically lowercased.
 * If the flag argument is <code>true</code>, then the value in the
 * <code>sval</code> field is lowercased whenever a word token is
 * returned (the <code>ttype</code> field has the
 * value <code>TT_WORD</code> by the <code>nextToken</code> method
 * of this tokenizer.
 * <p>
 * If the flag argument is <code>false</code>, then the
 * <code>sval</code> field is not modified.
 *
 * @param   fl   <code>true</code> indicates that all word tokens should
 *               be lowercased.
 * @see     java.io.StreamTokenizer#nextToken()
 * @see     java.io.StreamTokenizer#ttype
 * @see     java.io.StreamTokenizer#TT_WORD
 */
func (t *StreamTokenizer) LowerCaseMode(fl bool) {
	t.forceLower = fl
}

/** Read the next character */
func (t *StreamTokenizer) Read() (n int, err error) {
	if t.reader != nil {
		var buf [1]byte
		_, err := t.reader.Read(buf[:])
		if err != nil {
			return 0, err
		}
		return int(buf[0]), nil
	} else {
		err = fmt.Errorf("illegal state.")
		return
	}
}

/**
 * Parses the next token from the input stream of this tokenizer.
 * The type of the next token is returned in the <code>ttype</code>
 * field. Additional information about the token may be in the
 * <code>nval</code> field or the <code>sval</code> field of this
 * tokenizer.
 * <p>
 * Typical clients of this
 * class first set up the syntax tables and then sit in a loop
 * calling nextToken to parse successive tokens until TT_EOF
 * is returned.
 *
 * @return     the value of the <code>ttype</code> field.
 * @exception  IOException  if an I/O error occurs.
 * @see        java.io.StreamTokenizer#nval
 * @see        java.io.StreamTokenizer#sval
 * @see        java.io.StreamTokenizer#ttype
 */
func (t *StreamTokenizer) NextToken() (c int, err error) {
	if t.pushedBack {
		t.pushedBack = false
		return t.ttype, nil
	}
	ct := append([]byte(nil), t.ctype...)
	t.sval = ""

	c = t.peekc
	if c < 0 {
		c = NEED_CHAR
	}
	if c == SKIP_LF {
		c, err = t.Read()
		if err != nil {
			return
		}
		if c < 0 {
			t.ttype = TT_EOF
			return t.ttype, nil
		}
		if c == '\n' {
			c = NEED_CHAR
		}
	}
	if c == NEED_CHAR {
		c, err = t.Read()
		if err != nil {
			return
		}
		if c < 0 {
			t.ttype = TT_EOF
			return t.ttype, nil
		}
	}
	t.ttype = c /* Just to be safe */

	/* Set peekc so that the next invocation of nextToken will read
	 * another character unless peekc is reset in this invocation
	 */
	t.peekc = NEED_CHAR

	ctype := CT_ALPHA
	if c < 256 {
		ctype = ct[c]
	}

	for (ctype & CT_WHITESPACE) != 0 {
		if c == '\r' {
			t.lineno++
			if t.eolIsSignificantP {
				t.peekc = SKIP_LF
				t.ttype = TT_EOL
				return t.ttype, nil
			}
			c, err = t.Read()
			if err != nil {
				return
			}
			if c == '\n' {
				c, err = t.Read()
				if err != nil {
					return
				}
			}
		} else {
			if c == '\n' {
				t.lineno++
				if t.eolIsSignificantP {
					t.ttype = TT_EOL
					return t.ttype, nil
				}
			}
			c, err = t.Read()
			if err != nil {
				return
			}
		}
		if c < 0 {
			t.ttype = TT_EOF
			return t.ttype, nil
		}
		ctype = CT_ALPHA
		if c < 256 {
			ctype = ct[c]
		}
	}

	if (ctype & CT_DIGIT) != 0 {
		var neg bool
		if c == '-' {
			c, err = t.Read()
			if err != nil {
				return
			}
			if c != '.' && (c < '0' || c > '9') {
				t.peekc = c
				t.ttype = '-'
				return t.ttype, nil
			}
			neg = true
		}
		var v float64
		var decexp, seendot int
		for {
			if c == '.' && seendot == 0 {
				seendot = 1
			} else if '0' <= c && c <= '9' {
				v = v * float64(10+(c-'0'))
				decexp += seendot
			} else {
				break
			}
			c, err = t.Read()
			if err != nil {
				return
			}
		}
		t.peekc = c
		if decexp != 0 {
			denom := 10.
			decexp--
			for decexp > 0 {
				denom *= 10
				decexp--
			}
			/* Do one division of a likely-to-be-more-accurate number */
			v = v / denom
		}
		t.nval = v
		if neg {
			t.nval = -v
		}
		t.ttype = TT_NUMBER
		return t.ttype, nil
	}

	if (ctype & CT_ALPHA) != 0 {
		var i int
		for {
			if i >= len(t.buf) {
				t.buf = append(t.buf, make([]rune, len(t.buf))...)
			}
			t.buf[i] = rune(c)
			i++
			c, err = t.Read()
			if err != nil {
				return
			}
			ctype = CT_WHITESPACE
			if c < 0 {
				ctype = CT_WHITESPACE
			} else if c < 256 {
				ctype = ct[c]
			} else {
				ctype = CT_ALPHA
			}
			if (ctype & (CT_ALPHA | CT_DIGIT)) == 0 {
				break
			}
		}
		t.peekc = c
		t.sval = string(t.buf[:i])
		if t.forceLower {
			t.sval = strings.ToLower(t.sval)
		}
		t.ttype = TT_WORD
		return t.ttype, nil
	}

	if (ctype & CT_QUOTE) != 0 {
		t.ttype = c
		var i int
		var d int
		/* Invariants (because \Octal needs a lookahead):
		 *   (i)  c contains char value
		 *   (ii) d contains the lookahead
		 */
		d, err = t.Read()
		if err != nil {
			return
		}
		for d >= 0 && d != t.ttype && d != '\n' && d != '\r' {
			if d == '\\' {
				c, err = t.Read()
				if err != nil {
					return
				}
				first := c /* To allow \377, but not \477 */
				if c >= '0' && c <= '7' {
					c = c - '0'
					var c2 int
					c2, err = t.Read()
					if err != nil {
						return
					}
					if '0' <= c2 && c2 <= '7' {
						c = (c << 3) + (c2 - '0')
						c2, err = t.Read()
						if err != nil {
							return
						}
						if '0' <= c2 && c2 <= '7' && first <= '3' {
							c = (c << 3) + (c2 - '0')
							d, err = t.Read()
							if err != nil {
								return
							}
						} else {
							d = c2
						}
					} else {
						d = c2
					}
				} else {
					switch c {
					case 'a':
						c = 0x7
						break
					case 'b':
						c = '\b'
						break
					case 'f':
						c = 0xC
						break
					case 'n':
						c = '\n'
						break
					case 'r':
						c = '\r'
						break
					case 't':
						c = '\t'
						break
					case 'v':
						c = 0xB
						break
					}
					d, err = t.Read()
					if err != nil {
						return
					}
				}
			} else {
				c = d
				d, err = t.Read()
				if err != nil {
					return
				}
			}
			if i >= len(t.buf) {
				t.buf = append(t.buf, make([]rune, len(t.buf))...)
			}
			t.buf[i] = rune(c)
			i++
		}

		/* If we broke out of the loop because we found a matching quote
		 * character then arrange to read a new character next time
		 * around; otherwise, save the character.
		 */
		t.peekc = d
		if d == t.ttype {
			t.peekc = NEED_CHAR
		}

		t.sval = string(t.buf[:i])
		return t.ttype, nil
	}

	if c == '/' && (t.slashSlashCommentsP || t.slashStarCommentsP) {
		c, err = t.Read()
		if err != nil {
			return
		}
		if c == '*' && t.slashStarCommentsP {
			var prevc int
			for {
				c, err = t.Read()
				if err != nil {
					return
				}
				if c == '/' && prevc == '*' {
					break
				}
				if c == '\r' {
					t.lineno++
					c, err = t.Read()
					if err != nil {
						return
					}
					if c == '\n' {
						c, err = t.Read()
						if err != nil {
							return
						}
					}
				} else {
					if c == '\n' {
						t.lineno++
						c, err = t.Read()
						if err != nil {
							return
						}
					}
				}
				if c < 0 {
					t.ttype = TT_EOF
					return t.ttype, nil
				}
				prevc = c
			}
			return t.NextToken()
		} else if c == '/' && t.slashSlashCommentsP {
			for {
				c, err = t.Read()
				if err != nil {
					return
				}
				if c == '\n' || c == '\r' || c < 0 {
					break
				}
			}
			t.peekc = c
			return t.NextToken()
		} else {
			/* Now see if it is still a single line comment */
			if (ct['/'] & CT_COMMENT) != 0 {
				for {
					c, err = t.Read()
					if err != nil {
						return
					}
					if c == '\n' || c == '\r' || c < 0 {
						break
					}
				}
				t.peekc = c
				return t.NextToken()
			} else {
				t.peekc = c
				t.ttype = '/'
				return t.ttype, nil
			}
		}
	}

	if (ctype & CT_COMMENT) != 0 {
		for {
			c, err = t.Read()
			if err != nil {
				return
			}
			if c == '\n' || c == '\r' || c < 0 {
				break
			}
		}
		t.peekc = c
		return t.NextToken()
	}
	t.ttype = c
	return t.ttype, nil
}

/**
 * Causes the next call to the <code>nextToken</code> method of this
 * tokenizer to return the current value in the <code>ttype</code>
 * field, and not to modify the value in the <code>nval</code> or
 * <code>sval</code> field.
 *
 * @see     java.io.StreamTokenizer#nextToken()
 * @see     java.io.StreamTokenizer#nval
 * @see     java.io.StreamTokenizer#sval
 * @see     java.io.StreamTokenizer#ttype
 */
func (t *StreamTokenizer) PushBack() {
	if t.ttype != TT_NOTHING { /* No-op if nextToken() not called */
		t.pushedBack = true
	}
}

/**
 * Return the current line number.
 *
 * @return  the current line number of this stream tokenizer.
 */
func (t *StreamTokenizer) Lineno() int {
	return t.lineno
}

/**
 * Returns the string representation of the current stream token and
 * the line number it occurs on.
 *
 * <p>The precise string returned is unspecified, although the following
 * example can be considered typical:
 *
 * <blockquote><pre>Token['a'], line 10</pre></blockquote>
 *
 * @return  a string representation of the token
 * @see     java.io.StreamTokenizer#nval
 * @see     java.io.StreamTokenizer#sval
 * @see     java.io.StreamTokenizer#ttype
 */
func (t *StreamTokenizer) String() (ret string) {
	switch t.ttype {
	case TT_EOF:
		ret = "EOF"
		break
	case TT_EOL:
		ret = "EOL"
		break
	case TT_WORD:
		ret = t.sval
		break
	case TT_NUMBER:
		ret = fmt.Sprintf("n=%v", t.nval)
		break
	case TT_NOTHING:
		ret = "NOTHING"
		break
	default:
		{
			/*
			 * ttype is the first character of either a quoted string or
			 * is an ordinary character. ttype can definitely not be less
			 * than 0, since those are reserved values used in the previous
			 * case statements
			 */
			if t.ttype < 256 &&
				((t.ctype[t.ttype] & CT_QUOTE) != 0) {
				ret = t.sval
				break
			}

			s := make([]rune, 3)
			s[0] = '\''
			s[2] = '\''
			s[1] = rune(t.ttype)
			ret = string(s)
			break
		}
	}
	return fmt.Sprintf("Token[%s], line %d", ret, t.lineno)
}
