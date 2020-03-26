package chash

import (
    "hash/crc32"
    "sort"
    "strconv"
    "sync"
)

type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return "ConsistentError: " + e.s
}

func ConsistentError(text string) error {
    return &errorString{text}
}

type Circle []uint32

func (c Circle) Len() int {
    return len(c)
}

func (c Circle) Less(i, j int) bool {
    return c[i] < c[j]
}

func (c Circle) Swap(i, j int) {
    c[i], c[j] = c[j], c[i]
}

type Hash func(date []byte) uint32

type Consistent struct {
    hash         Hash
    circle       Circle
    virtualNodes int
    virtualMap   map[uint32]string
    members      map[string]bool
    sync.RWMutex
}

func NewCHash(vnode int) *Consistent {
	if vnode<=0{
		panic("vnode<=0")
	}
    return &Consistent{
	    hash:         crc32.ChecksumIEEE,
	    circle:       Circle{},
	    virtualNodes: vnode,
	    virtualMap:   make(map[uint32]string),
	    members:      make(map[string]bool),
    }
}

func (c *Consistent) eltKey(key string, idx int) string {
    return key + "|" + strconv.Itoa(idx)
}

func (c *Consistent) updateCricle() {
    c.circle = Circle{}
    for k := range c.virtualMap {
	    c.circle = append(c.circle, k)
    }
    sort.Sort(c.circle)
}

func (c *Consistent) Members() []string {
    c.RLock()
    defer c.RUnlock()

    m := make([]string, len(c.members))

    var i = 0
    for k := range c.members {
	    m[i] = k
	    i++
    }

    return m
}

func (c *Consistent) Get(key string) string {
    hashKey := c.hash([]byte(key))
    c.RLock()
    defer c.RUnlock()

    i := c.search(hashKey)

    return c.virtualMap[c.circle[i]]
}

func (c *Consistent) search(key uint32) int {
    f := func(x int) bool {
	    return c.circle[x] >= key
    }

    i := sort.Search(len(c.circle), f)
    i = i - 1
    if i < 0 {
	    i = len(c.circle) - 1
    }
    return i
}

func (c *Consistent) ForceSet(keys ...string) {
    mems := c.Members()
    for _, elt := range mems {
	    var found = false

    FOUNDLOOP:
	    for _, k := range keys {
		    if k == elt {
			    found = true
			    break FOUNDLOOP
		    }
	    }
	    if !found {
		    c.Remove(elt)
	    }
    }

    for _, k := range keys {
	    c.RLock()
	    _, ok := c.members[k]
	    c.RUnlock()

	    if !ok {
		    c.Add(k)
	    }
    }
}

func (c *Consistent) Add(elt string) {
    c.Lock()
    defer c.Unlock()

    if _, ok := c.members[elt]; ok {
	    return
    }

    c.members[elt] = true

    for idx := 0; idx < c.virtualNodes; idx++ {
	    c.virtualMap[c.hash([]byte(c.eltKey(elt, idx)))] = elt
    }

    c.updateCricle()
}

func (c *Consistent) Remove(elt string) {
    c.Lock()
    defer c.Unlock()

    if _, ok := c.members[elt]; !ok {
	    return
    }

    delete(c.members, elt)

    for idx := 0; idx < c.virtualNodes; idx++ {
	    delete(c.virtualMap, c.hash([]byte(c.eltKey(elt, idx))))
    }

    c.updateCricle()
}
