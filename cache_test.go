package lilcache

import (
  "fmt"
  "testing"
)

func TestBasic(t *testing.T) {
  c := New(10)

  // not found
  v, tm := c.Get("foo")
  if v != nil || !tm.IsZero() {
    t.Errorf("expected nil v, zero t, got %v and %v", v, tm)
  }

  c.Put("foo", "val")
  v, tm = c.Get("foo")
  if v != "val" || tm.IsZero() {
    t.Errorf("expected v of \"val\" and non-zero t, got %v and %v", v, tm)
  }
}

func TestCap(t *testing.T) {
  c := New(3)

  for i := 0; i < 3; i++ {
    k := fmt.Sprintf("%d", i)
    c.Put(k, k)
  }

  for i := 0; i < 3; i++ {
    k := fmt.Sprintf("%d", i)
    if v, tm := c.Get(k); k != v || tm.IsZero() {
      t.Errorf("expected key (%s) with val (%s), got %v t=%v", k, k, v, tm)
    }
  }

  // invalidate "0"
  c.Put("3", "3")
  if v, tm := c.Get("0"); v != nil || !tm.IsZero() {
    t.Errorf("expected eviction of \"0\", got %v t=%v", v, tm)
  }

  for i := 1; i < 4; i++ {
    k := fmt.Sprintf("%d", i)
    if v, tm := c.Get(k); k != v || tm.IsZero() {
      t.Errorf("expected key (%s) with (%s), got %v t=%v", k, k, v, tm)
    }
  }
}

func TestDelete(t *testing.T) {
  c := New(1)

  c.Put("1", "1")

  if v, tm := c.Delete("1"); v != "1" || tm.IsZero() {
    t.Errorf("expected \"1\" and non-zero time, got %v and %v", v, tm)
  }

  if v, tm := c.Delete("1"); v != nil || !tm.IsZero() {
    t.Errorf("expected nil and zero time, got %v and %v", v, tm)
  }
}
