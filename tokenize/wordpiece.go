package tokenize

import (
	"fmt"

	"github.com/lfpython/gobert/tokenize/vocab"
)

// DefaultMaxWordChars is the max length of a token for it to be tokenized, otherwise marked as unknown
const DefaultMaxWordChars = 200

// DefaultUnknownToken is the token used to signify an unkown token
const DefaultUnknownToken = "[UNK]"

// Wordpiece is a tokenizer that breaks tokens into subword units based on a supplied vocabulary
// https://arxiv.org/pdf/1609.08144.pdf Section 4.1 for details
type Wordpiece struct {
	vocab        vocab.Dict
	maxWordChars int
	unknownToken string
}

// NewWordpiece returns a WordpieceTokenizer with the default settings.
// Generally should be used in a FullTokenizer
func NewWordpiece(voc vocab.Dict) Wordpiece {
	return Wordpiece{
		vocab:        voc,
		maxWordChars: DefaultMaxWordChars,
		unknownToken: DefaultUnknownToken,
	}
}

// Tokenize will segment the text into subword tokens from the supplied vocabulary
// NOTE: This implementation does not EXACTLY match the ref-impl and behaves slightly differently
// See https://github.com/google-research/bert/issues/763
func (wp Wordpiece) Tokenize(text string) []string {
	// TODO: determine if utf8 conversion is necessary, per python impl
	// text = convert_to_unicode(text)
	var toks []string
	for _, tok := range tokenizeWhitespace(text) {
		if len(tok) > wp.maxWordChars {
			toks = append(toks, wp.unknownToken)
			continue
		}
		var tks []string
		tk := tok
		for len(tk) > 0 && tk != "##" {
			sub := wp.vocab.LongestSubstring(tk)
			if sub == "" {
				tks = []string{tok}
				break
			}
			tks = append(tks, sub)
			tk = fmt.Sprintf("##%s", tk[len(sub):])
		}
		toks = append(toks, tks...)
	}
	return toks
}

// SetMaxWordChars will set the max chars for a word to be tokenized,
// generally this should be congfigured through the FullTokenizer
func (wp Wordpiece) SetMaxWordChars(c int) {
	wp.maxWordChars = c
}

// SetUnknownToken will set the , generally this should be congfigured through the FullTokenizer
func (wp Wordpiece) SetUnknownToken(tok string) {
	wp.unknownToken = tok
}
