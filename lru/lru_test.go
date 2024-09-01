package lru

import (
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	c := New(0, nil)
	c.AddOrSet("k1", String("v1"))
	if v, ok := c.Get("k1"); !ok && v.(String) != "v1" {
		t.Fatalf("get k1 err")
	}
	if _, ok := c.Get("k2"); ok {
		t.Fatalf("get k2 err")
	}
}

func TestAddOrGet(t *testing.T) {
	c := New(12, nil)
	c.AddOrSet("k1", String("v1"))
	c.AddOrSet("k2", String("v2"))
	c.AddOrSet("k3", String("v3"))
	c.AddOrSet("k4", String("v4"))
	if _, ok := c.Get("k1"); ok {
		t.Fatalf("k1 err")
	}
	c.AddOrSet("k2", String("2v"))
	if v, ok := c.Get("k2"); !ok && v.(String) != "2v" {
		t.Fatalf("k2 err")
	}
	c.AddOrSet("k5", String("v5"))
	if _, ok := c.Get("k3"); ok {
		t.Fatalf("k3 err")
	}
	if c.Len() != 3 {
		t.Fatalf("len err: %v", c.Len())
	}
}
