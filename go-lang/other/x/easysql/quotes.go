package easysql

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func isSplitRune(r rune) bool {
	switch r {
	case ' ', '(', ')', '<', '>', '!', '=', '?':
		return true
	}
	return false
}

func isSplitToken(s string) bool {
	ns := strings.ToUpper(s)
	if ns == "AND" || ns == "OR" {
		return true
	}
	if len(s) == 1 {
		return isSplitRune(rune(s[0]))
	}
	return false
}

func scanToken(rd *strings.Reader) (token string, err error) {
	r, _, err := rd.ReadRune()
	if err != nil {
		return "", err
	}
	if isSplitRune(r) {
		return string(r), nil
	}

	var tok []rune
	for {
		tok = append(tok, r)
		r, _, err = rd.ReadRune()
		if err != nil {
			if len(tok) > 0 {
				return string(tok), nil
			}
			return "", err
		}
		if isSplitRune(r) {
			rd.UnreadRune()
			return string(tok), nil
		}
	}
}

func sqlQuotes(s string) string {
	var buf bytes.Buffer
	r := strings.NewReader(s)
	for {
		tok, err := scanToken(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "scan token: %v\n", err)
			break
		}
		if isSplitToken(tok) {
			fmt.Fprintf(&buf, tok)
		} else {
			fmt.Fprintf(&buf, "`%s`", tok)
		}
	}
	return buf.String()
}
