package parser

type lexer struct {
	src     []byte
	ch      byte
	offset  int
	pos     int
	nextPos int
}

func newLexer(src []byte) *lexer {
	l := &lexer{src: src}
	l.next()
	return l
}

func (l *lexer) scan() (int, token, string) {
	tok := ILLEGAL
	pos := l.pos
	val := ""

	if l.ch == 0 {
		return l.pos, EOL, ""
	}

	// return the number of spaces.
	// We need this to figure out if the transaction is debit or credit.
	if isSpace(l.ch) {
		start := l.offset - 1
		for l.ch == ' ' {
			l.next()
		}
		return l.pos, PAD, string(l.src[start : l.offset-1])
	}

	if isSymbol(l.ch) {
		l.next()
		return pos, SYMBOL, string(l.ch)
	}

	if isAlpha(l.ch) {
		start := l.offset - 1
		for isAlpha(l.ch) || isDigit(l.ch) || isSymbol(l.ch) {
			l.next()
		}
		return l.pos, WORD, string(l.src[start : l.offset-1])
	}

	if isDigit(l.ch) {
		start := l.offset - 1
		// to be a date, we need 2 digits + dot + 2 digits
		countPrefixDigits := 1
		countDot := 0
		countSuffixDigits := 0
		for isDigit(l.ch) || isDot(l.ch) || isComma(l.ch) {
			if isDot(l.ch) {
				countDot += 1
				l.next()
				continue
			}
			switch countDot {
			case 0:
				countPrefixDigits += 1
			default:
				countSuffixDigits += 1
			}

			l.next()
		}

		if countPrefixDigits == 2 && countDot == 1 && countSuffixDigits == 2 {
			return l.pos, SHORT_DATE, string(l.src[start : l.offset-1])
		}

		return l.pos, NUMBER, string(l.src[start : l.offset-1])
	}

	return pos, tok, val
}

// Load the next character into l.ch (or 0 on end of input) and update line position.
func (l *lexer) next() {
	l.pos = l.nextPos
	if l.offset >= len(l.src) {
		// For last character, move offset 1 past the end as it
		// simplifies offset calculations in NAME and NUMBER
		if l.ch != 0 {
			l.ch = 0
			l.offset++
			l.nextPos++
		}
		return
	}
	ch := l.src[l.offset]
	l.ch = ch
	l.nextPos++
	l.offset++
}

// IsAlpha checks if the given byte is an alphabetic character (a-z, A-Z).
func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// IsDot checks if the given byte is a dot character.
func isDot(ch byte) bool {
	return ch == '.'
}

// IsComma checks if the given byte is a comma character.
func isComma(ch byte) bool {
	return ch == ','
}

// IsSymbol checks if the given byte is one of the recognized symbol characters.
func isSymbol(ch byte) bool {
	return ch == '*' ||
		ch == '-' ||
		ch == '/' ||
		ch == ':' ||
		ch == '(' ||
		ch == ')' ||
		ch == '°'
}

// IsDigit checks if the given byte is a numeric digit (0-9).
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// IsSpace checks if the given byte is a space character.
func isSpace(ch byte) bool {
	return ch == ' '
}
