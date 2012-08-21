package lilcache

import (
  "fmt"
  "testing"
  "time"
)

func expectThis(t *testing.T, e, a interface{}, tm time.Time) {
  if e != a || tm.IsZero() {
    t.Errorf("expected \"%v\" and non-zero t, got %v and %v", e, a, tm)
  }
}

func expectEmpty(t *testing.T, a interface{}, tm time.Time) {
  if a != nil || !tm.IsZero() {
    t.Errorf("expected nil and zero t, got %v and %v", a, tm)
  }
}

func TestBasic(t *testing.T) {
  c := New(10)

  // not found
  v, tm := c.Get("foo")
  expectEmpty(t, v, tm)

  c.Put("foo", "val")
  v, tm = c.Get("foo")
  expectThis(t, "val", v, tm)
}

func TestCap(t *testing.T) {
  c := New(3)

  for i := 0; i < 3; i++ {
    k := fmt.Sprintf("%d", i)
    c.Put(k, k)
  }

  for i := 0; i < 3; i++ {
    k := fmt.Sprintf("%d", i)
    v, tm := c.Get(k)
    expectThis(t, k, v, tm)
  }

  // invalidate "0"
  c.Put("3", "3")
  v, tm := c.Get("0")
  expectEmpty(t, v, tm)

  for i := 1; i < 4; i++ {
    k := fmt.Sprintf("%d", i)
    v, tm := c.Get(k)
    expectThis(t, k, v, tm)
  }
}

func TestDelete(t *testing.T) {
  c := New(1)

  c.Put("1", "1")

  v, tm := c.Delete("1")
  expectThis(t, "1", v, tm)

  v, tm = c.Delete("1")
  expectEmpty(t, v, tm)
}

func TestInterleaved(t *testing.T) {
  c := New(4)

  c.Put("a", "a")
  c.Put("b", "b")

  v, tm := c.Delete("a")
  expectThis(t, "a", v, tm)

  c.Put("c", "c")
  v, tm = c.Delete("b")
  expectThis(t, "b", v, tm)

  v, tm = c.Delete("c")
  expectThis(t, "c", v, tm)

  v, tm = c.Get("a")
  expectEmpty(t, v, tm)
  v, tm = c.Get("b")
  expectEmpty(t, v, tm)
  v, tm = c.Get("c")
  expectEmpty(t, v, tm)
}
