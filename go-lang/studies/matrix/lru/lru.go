package lru

import "container/list"

type Key interface{}

type Cache struct {
	maxEntries int

	onEvicted func(key Key, value interface{})

	ll    *list.List
	cache map[Key]*list.Element
}

type entry struct {
	key   Key
	value interface{}
}

func New(maxEntries int, onEvicted func(key Key, value interface{})) *Cache {
	return new(Cache).Init(maxEntries, onEvicted)
}

func (c *Cache) Init(maxEntries int, onEvicted func(key Key, value interface{})) *Cache {
	c.maxEntries = maxEntries
	c.onEvicted = onEvicted
	c.ll = list.New()
	c.cache = make(map[Key]*list.Element)
	return c
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) Add(key Key, value interface{}) {
	if le, ok := c.cache[key]; ok {
		c.ll.MoveToFront(le)
		le.Value.(*entry).value = value
		return
	}
	le := c.ll.PushFront(&entry{key, value})
	c.cache[key] = le
	if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key Key) (interface{}, bool) {
	if le, hit := c.cache[key]; hit {
		c.ll.MoveToFront(le)
		return le.Value.(*entry).value, true
	}
	return nil, false
}

func (c *Cache) Remove(key Key) {
	if le, hit := c.cache[key]; hit {
		c.removeElement(le)
	}
}

func (c *Cache) RemoveOldest() {
	le := c.ll.Back()
	if le != nil {
		c.removeElement(le)
	}
}

func (c *Cache) removeElement(le *list.Element) {
	c.ll.Remove(le)
	kv := le.Value.(*entry)
	delete(c.cache, kv.key)
	if c.onEvicted != nil {
		c.onEvicted(kv.key, kv.value)
	}
}
