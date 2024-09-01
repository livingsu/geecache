package geecache

import (
	"errors"
	"testing"
)

func TestGet(t *testing.T) {
	kvs := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}
	visitCount := make(map[string]int)

	g := NewGroup("test", 100, GetterFunc(func(key string) ([]byte, error) {
		if v, ok := kvs[key]; ok {
			visitCount[key] += 1
			return []byte(v), nil
		}
		return nil, errors.New("key not found")
	}))
	if v, err := g.Get("k1"); err != nil || v.String() != "v1" || visitCount["k1"] != 1 {
		t.Fatalf("get err")
	}
	if v, err := g.Get("k1"); err != nil || v.String() != "v1" || visitCount["k1"] != 1 {
		t.Fatalf("get err")
	}
}
