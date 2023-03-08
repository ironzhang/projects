package iciba

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const key = "9703E5C2F64060501DB26311D45CF2CA"

var ErrNotFound = errors.New("not found")

type Exchange struct {
	WordPL    []string `json:"word_pl,omitempty"`
	WordPAST  []string `json:"word_past,omitempty"`
	WordDONE  []string `json:"word_done,omitempty"`
	WordING   []string `json:"word_ing,omitempty"`
	WordTHIRD []string `json:"word_third,omitempty"`
	WordER    []string `json:"word_er,omitempty"`
	WordEST   []string `json:"word_est,omitempty"`
}

type Part struct {
	Part  string   `json:"part,omitempty"`
	Means []string `json:"means,omitempty"`
}

type Symbol struct {
	PhEN      string `json:"ph_en,omitempty"`
	PhAM      string `json:"ph_am,omitempty"`
	PhOther   string `json:"ph_other,omitempty"`
	PhEN_MP3  string `json:"ph_en_mp3,omitempty"`
	PhAM_MP3  string `json:"ph_am_mp3,omitempty"`
	PhTTS_MP3 string `json:"ph_tts_mp3,omitempty"`
	Parts     []Part `json:"parts,omitempty"`
}

type Result struct {
	WordName string `json:"word_name,omitempty"`
	//Exchange Exchange `json:"exchange"`
	Symbols []Symbol `json:"symbols,omitempty"`
}

func Lookup(word string) (res Result, err error) {
	const URL = "http://dict-co.iciba.com/api/dictionary.php"

	params := make(url.Values)
	params.Add("key", key)
	params.Add("type", "json")
	params.Add("w", strings.ToLower(word))
	resp, err := http.Get(URL + "?" + params.Encode())
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, err
	}
	if res.WordName == "" || len(res.Symbols) <= 0 {
		return res, ErrNotFound
	}
	return res, nil
}
