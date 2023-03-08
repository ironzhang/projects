package pattern

import (
	"fmt"
	"testing"
)

func TestNode(t *testing.T) {
	items := []struct {
		pats  []string
		value interface{}
	}{
		{
			pats:  []string{"account", "v1", "login", ":user_id"},
			value: 0,
		},
		{
			pats:  []string{"account", "v1", ":user_id", "login"},
			value: 0,
		},
		{
			pats:  []string{"account", "v1", ":uid", "logout"},
			value: 0,
		},
	}

	n := node{name: "POST", children: make(map[string]*node)}
	for _, item := range items {
		if err := n.Add(item.pats, item.value); err != nil {
			t.Fatal(err)
		}
	}
	fmt.Println(n.String())
}
