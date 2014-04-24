package mem

import (
		"github.com/lexlapax/graveldb/core"
		"sync"
		"sort"
)

type edgeSet struct {
	atommap map[string]core.Edge
	sync.RWMutex
}

func newEdgeSet() *edgeSet {
	set := new(edgeSet)
	set.atommap = make(map[string]core.Edge)
	return set
}

func (set *edgeSet) add(edge core.Edge) {
	if edge == nil || edge.Id() == nil { return }
	id := string(edge.Id()[:])
	set.Lock()
	defer set.Unlock()
	if _, ok := set.atommap[id]; ok {
		return
	} else {
		set.atommap[id] = edge
	}
	return
}

func (set *edgeSet) del(edge core.Edge) {
	if edge == nil || edge.Id() == nil { return }
	id := string(edge.Id()[:])
	set.Lock()
	defer set.Unlock()
	delete(set.atommap, id)
	return
}

func (set *edgeSet) exists(edge core.Edge) bool {
	if edge == nil || edge.Id() == nil { return false}
	if _, ok := set.atommap[string(edge.Id()[:])]; ok {
		return true
	}
	return false
}

func (set *edgeSet) members() []core.Edge {
	atoms := []core.Edge{}
	keys := []string{}
	set.RLock()
	defer set.RUnlock()
	for k, _ := range set.atommap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		atoms = append(atoms, set.atommap[k])
	}
	return atoms
}

func (set *edgeSet) count() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.atommap)
}
