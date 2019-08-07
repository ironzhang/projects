package pattern

import (
	"bytes"
	"fmt"
	"io"
)

func patternName(name string) (string, bool) {
	if name[0] == ':' {
		return name[1:], true
	}
	return name, false
}

type node struct {
	name     string
	subnodes map[string]*node
	patnodes map[string]*node
	value    interface{}
}

func (n *node) Add(pats []string, val interface{}) error {
	var c *node
	var err error

	t := n
	for _, p := range pats {
		if c, err = t.addChild(p); err != nil {
			return err
		}
		t = c
	}
	t.value = val

	return nil
}

func (n *node) Get(pats []string) (map[string]string, error) {
	t := n
	values := make(map[string]string)
	for _, p := range pats {
		c, err := t.getChild(p)
		if err != nil {
			return nil, err
		}
		if c.pattern {
			values[c.name] = p
		}
		t = c
	}
	return values, nil
}

func (n *node) addChild(name string) (*node, error) {
	if name, pattern := patternName(name); pattern {
		return n.addPatternNode(name)
	} else {
		return n.addChildNode(name)
	}
}

func (n *node) addChildNode(name string) (*node, error) {
	if c, ok := n.children[name]; ok {
		return c, nil
	}
	c := &node{
		pattern:  false,
		name:     name,
		children: make(map[string]*node),
	}
	n.children[name] = c
	return c, nil
}

func (n *node) addPatternNode(name string) (*node, error) {
	if n.patnode != nil {
		if n.patnode.name == name {
			return n.patnode, nil
		} else {
			return nil, fmt.Errorf("pattern conflict: %q != %q", n.patnode.name, name)
		}
	}
	c := &node{
		pattern:  true,
		name:     name,
		children: make(map[string]*node),
	}
	n.patnode = c
	return c, nil
}

func (n *node) getChild(name string) (*node, error) {
	if c, ok := n.children[name]; ok {
		return c, nil
	}
	if n.patnode != nil {
		return n.patnode, nil
	}
	return nil, fmt.Errorf("%q node not found", name)
}

func (n *node) String() string {
	var buf bytes.Buffer
	n.Print(&buf, 0)
	return buf.String()
}

func (n *node) Print(w io.Writer, d int) {
	fmt.Fprint(w, "|")
	for i := 0; i < d; i++ {
		fmt.Fprint(w, "-")
	}
	if n.pattern {
		fmt.Fprint(w, ":")
	}
	fmt.Fprint(w, n.name)
	if n.value != nil {
		fmt.Fprint(w, "\t", n.value)
	}
	fmt.Fprintln(w)

	d++
	for _, c := range n.children {
		c.Print(w, d)
	}
	if n.patnode != nil {
		n.patnode.Print(w, d)
	}
}
