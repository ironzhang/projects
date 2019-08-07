package gosdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ironzhang/wordbook/cores/types"
)

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{url: url}
}

func (c *Client) ListWords(offset, limit int) ([]types.Word, error) {
	path := fmt.Sprintf("/api/words?offset=%d&limit=%d", offset, limit)
	resp, err := http.Get(c.url + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status is no ok")
	}

	var words []types.Word
	err = json.NewDecoder(resp.Body).Decode(&words)
	if err != nil {
		return nil, err
	}
	return words, nil
}

func (c *Client) LookupWord(word string) (types.Word, error) {
	path := fmt.Sprintf("/api/words/%s", word)
	resp, err := http.Get(c.url + path)
	if err != nil {
		return types.Word{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Word{}, errors.New("status is not ok")
	}

	var w types.Word
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return types.Word{}, err
	}
	return w, nil
}

func (c *Client) AdjustWordPriority(word string, n int) error {
	path := fmt.Sprintf("/api/words/%s/p?n=%d", word, n)
	req, err := http.NewRequest("PUT", c.url+path, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("status is not ok")
	}
	return nil
}
