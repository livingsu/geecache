package lru

import (
	"container/list"
)

const defaultMaxBytes = 10 * 1024 * 1024 // 10MB

type Cache struct {
	l         *list.List
	m         map[string]*list.Element
	maxBytes  int64
	curBytes  int64
	OnEvicted func(key string, value any)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(key string, value any)) *Cache {
	if maxBytes <= 0 {
		maxBytes = defaultMaxBytes
	}
	return &Cache{
		l:         list.New(),
		m:         make(map[string]*list.Element),
		maxBytes:  maxBytes,
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value any, ok bool) {
	if e, ok := c.m[key]; ok {
		c.l.MoveToFront(e)
		value = e.Value.(*entry).value
		return value, true
	}
	return
}

func (c *Cache) removeOldest() {
	oldest := c.l.Back()
	if oldest != nil {
		e := oldest.Value.(*entry)
		c.l.Remove(oldest)
		delete(c.m, e.key)
		c.curBytes = c.curBytes - int64(len(e.key)) - int64(e.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(e.key, e.value)
		}
	}
}

func (c *Cache) AddOrSet(key string, value Value) {
	if e, ok := c.m[key]; ok {
		c.l.MoveToFront(e)
		kv := e.Value.(*entry)
		kv.value = value
		c.curBytes += int64(value.Len()) - int64(kv.value.Len())
	} else {
		c.m[key] = c.l.PushFront(&entry{key, value})
		c.curBytes += int64(len(key)) + int64(value.Len())
	}
	for c.curBytes > c.maxBytes {
		c.removeOldest()
	}
}

func (c *Cache) Len() int {
	return c.l.Len()
}
