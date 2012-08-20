package lilcache

import (
  "sync"
  "time"
)

type Cache struct {
  l          sync.Mutex
  m          map[string]*entry
  head, tail *entry
  size       int
  capacity   int
}

type entry struct {
  key        string
  val        interface{}
  next, prev *entry
  t          time.Time
}

func head(c *Cache, e *entry) {
  e.t = time.Now()

  if c.head == e {
    return
  }

  // take e out of its current predicament
  if n := e.next; n != nil {
    n.prev = e.prev
  }

  if p := e.prev; p != nil {
    p.next = e.next
  }

  if h := c.head; h != nil {
    h := c.head
    h.prev = e
    e.next = h
    e.prev = nil
  } else {
    c.head = e
    c.tail = e
  }
}

func tail(c *Cache) {
  for c.size > c.capacity {
    t := c.tail
    c.tail = t.prev
    delete(c.m, t.key)
    c.size -= 1
  }
}

func New(cap int) *Cache {
  return &Cache{
    m:        make(map[string]*entry),
    capacity: cap,
  }
}

// An empty time on return means non-existence.
func (c *Cache) Get(key string) (interface{}, time.Time) {
  c.l.Lock()
  defer c.l.Unlock()

  e, ok := c.m[key]
  if !ok {
    return nil, time.Time{}
  }

  head(c, e)
  return e.val, e.t
}

func (c *Cache) Put(key string, val interface{}) {
  c.l.Lock()
  defer c.l.Unlock()

  e, ok := c.m[key]
  if ok {
    e.val = val
    head(c, e)
    return
  }

  e = &entry{key: key, val: val}
  c.m[key] = e
  c.size++
  head(c, e)
  tail(c)
}
