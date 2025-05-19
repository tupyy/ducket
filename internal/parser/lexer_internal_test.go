package parser

import (
	"strings"
	"testing"
)

func TestTokens(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "20.20 20.2020",
			output: "short_date number EOL",
		},
		// {
		// 	input:  "1 200,00",
		// 	output: "number EOL",
		// },
		{
			input:  "20/10",
			output: "number symbol number EOL",
		},
		{
			input:  "200.20 0002 20,20 1 200,00",
			output: "number number number number EOL",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := newLexer([]byte(test.input))

			tokens := []string{}
			for {
				_, tok, _ := l.scan()
				tokens = append(tokens, tok.String())
				if tok == EOL {
					break
				}

			}

			output := strings.Join(tokens, " ")
			if strings.TrimSpace(output) != test.output {
				t.Errorf("expected %q, got %q", test.output, output)
			}
		})
	}
}
