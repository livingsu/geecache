package geecache

import (
	"errors"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Print("[geecache] get hit\n")
		return v, nil
	}
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	b, err := g.getter.Get(key)
	if err != nil {
		log.Printf("[geecache] getter err: %s\n", err)
		return ByteView{}, err
	}
	log.Printf("[geecache] getter get key: %s\n", key)
	v := ByteView{cloneBytes(b)}
	g.mainCache.add(key, v)
	return v, nil
}
