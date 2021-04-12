package bookserver

import (
	"database/sql"

	"github.com/ironzhang/wordbook/cores/storage"
	"github.com/ironzhang/wordbook/cores/types"
)

type Server struct {
	lookup     func(word string) (types.Word, error)
	words      *storage.WordTable
	priorities *storage.PriorityTable
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		lookup:     lookupFromICIBA,
		words:      storage.NewWordTable(db),
		priorities: storage.NewPriorityTable(db),
	}
}

func (s *Server) AdjustWordPriority(word string, n int) error {
	p, ok, err := s.priorities.Get(word)
	if err != nil {
		return err
	}
	if ok {
		return s.priorities.Set(word, p+n)
	}
	return s.priorities.Set(word, 10+n)
}

func (s *Server) LookupWord(word string) (w types.Word, ok bool, err error) {
	w, ok, err = s.words.Get(word)
	if err != nil {
		return w, false, err
	}
	s.AdjustWordPriority(word, 1)
	if ok {
		return w, true, nil
	}
	if w, err = s.lookup(word); err != nil {
		return w, false, err
	}
	s.words.Set(w)
	return w, true, nil
}

func (s *Server) ListWords(offset, limit int, order bool) ([]types.Word, error) {
	if order {
		return s.words.ListByPriority(offset, limit)
	}
	return s.words.List(offset, limit)
}
