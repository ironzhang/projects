package main

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/alecthomas/template"
	"github.com/ironzhang/wordbook/cores/types"
	"github.com/ironzhang/wordbook/sdks/gosdk"
	"github.com/ironzhang/x-pearls/log"
)

type site struct {
	client *gosdk.Client
}

func (s *site) Home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		word, err := s.client.LookupWord(r.Form.Get("lookup-text"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t := template.Must(template.ParseFiles("./html/index.template.html"))
		if err = t.Execute(w, word); err != nil {
			log.Errorf("Execute: %v", err)
		}
	default:
		data, err := ioutil.ReadFile("./html/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

type WordsPageData struct {
	Words []types.Word
	Page  int
}

func (s *site) ListWords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		page, err := strconv.Atoi(r.Form.Get("page-text"))
		if err != nil {
			page = 0
		}
		words, err := s.client.ListWords(10*page, 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t := template.Must(template.ParseFiles("./html/words.template.html"))
		if err = t.Execute(w, WordsPageData{Words: words, Page: page + 1}); err != nil {
			log.Errorf("Execute: %v", err)
		}
	default:
		words, err := s.client.ListWords(0, 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t := template.Must(template.ParseFiles("./html/words.template.html"))
		if err = t.Execute(w, WordsPageData{Words: words, Page: 1}); err != nil {
			log.Errorf("Execute: %v", err)
		}
	}
}

func (s *site) Priority(w http.ResponseWriter, r *http.Request) {
	words, err := s.client.ListWords(0, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t := template.Must(template.ParseFiles("./html/priority.template.html"))
	if err = t.Execute(w, WordsPageData{Words: words, Page: 1}); err != nil {
		log.Errorf("Execute: %v", err)
	}
}

func (s *site) Adjust(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	word := params.Get("word")
	action := params.Get("action")

	n := 0
	switch action {
	case "inc":
		n = 1
	case "dec":
		n = -1
	}

	if err := s.client.AdjustWordPriority(word, n); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	s := site{client: gosdk.NewClient("http://localhost:8080")}
	http.HandleFunc("/", s.Home)
	http.HandleFunc("/words", s.ListWords)
	http.HandleFunc("/priority", s.Priority)
	http.HandleFunc("/adjust", s.Adjust)
	http.ListenAndServe(":8081", nil)
}
