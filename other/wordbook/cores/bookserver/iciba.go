package bookserver

import (
	"github.com/ironzhang/wordbook/cores/iciba"
	"github.com/ironzhang/wordbook/cores/types"
)

func wrap(w types.Word) types.Word {
	if w.Sound.EN == "" {
		w.Sound.EN = w.Sound.TTS
	}
	if w.Sound.AM == "" {
		w.Sound.AM = w.Sound.TTS
	}
	return w
}

func lookupFromICIBA(word string) (types.Word, error) {
	res, err := iciba.Lookup(word)
	if err != nil {
		return types.Word{}, err
	}
	w := types.Word{Word: word}
	if len(res.Symbols) > 0 {
		w.PhEN = res.Symbols[0].PhEN
		w.PhAM = res.Symbols[0].PhAM
		w.PhOther = res.Symbols[0].PhOther
		w.Sound.EN = res.Symbols[0].PhEN_MP3
		w.Sound.AM = res.Symbols[0].PhAM_MP3
		w.Sound.TTS = res.Symbols[0].PhTTS_MP3
		for _, p := range res.Symbols[0].Parts {
			w.Parts = append(w.Parts, types.Part{Part: p.Part, Means: p.Means})
		}
	}
	return wrap(w), nil
}
