package parser

type token int

const (
	ILLEGAL token = iota
	WORD
	DATE
	SHORT_DATE
	PAD // number of spaces
	NUMBER
	SYMBOL
	EOL
)

var tokenNames = map[token]string{
	ILLEGAL:    "illegal",
	EOL:        "EOL",
	PAD:        "padding",
	WORD:       "word",
	NUMBER:     "number",
	DATE:       "date",
	SHORT_DATE: "short_date",
	SYMBOL:     "symbol",
}

func (t token) String() string {
	return tokenNames[t]
}
