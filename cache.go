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

func extract(e *entry) {
  if n := e.next; n != nil {
    n.prev = e.prev
  }

  if p := e.prev; p != nil {
    p.next = e.next
  }
}

func head(c *Cache, e *entry) {
  if c.head == e {
    return
  }

  extract(e)

  if h := c.head; h != nil {
    h.prev = e
    e.next = h
    e.prev = nil
  } else {
    c.tail = e
  }
  c.head = e
}

func tail(c *Cache) {
  for c.size > c.capacity {
    delete(c.m, c.tail.key)
    c.tail = c.tail.prev
    c.tail.next = nil
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

func (c *Cache) Delete(key string) (interface{}, time.Time) {
  c.l.Lock()
  defer c.l.Unlock()

  e, ok := c.m[key]
  if !ok {
    return nil, time.Time{}
  }

  delete(c.m, key)

  extract(e)

  if c.head == e {
    c.head = e.next
  }

  if c.tail == e {
    c.tail = e.prev
  }

  c.size--
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

  e = &entry{key: key, val: val, t: time.Now()}
  c.m[key] = e
  c.size++
  head(c, e)
  tail(c)
}

// func (c *Cache) Dump() {
//   fmt.Printf("Size: %d\n", c.size)
//   fmt.Println("Map")
//   for k, v := range c.m {
//     fmt.Printf("%v -> %v\n", k, v)
//   }

//   fmt.Println("Lst")
//   for h := c.head; h != nil; h = h.next {
//     fmt.Printf("%v -> %v\n", h.key, h.val)
//   }
//   fmt.Println()
// }
