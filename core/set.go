package core

import (
		"sync"
		"sort"
)

type EdgeSet struct {
	atommap map[string]Edge
	sync.RWMutex
}

func NewEdgeSet() *EdgeSet {
	set := new(EdgeSet)
	set.atommap = make(map[string]Edge)
	return set
}

func (set *EdgeSet) Add(edge Edge) {
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

func (set *EdgeSet) Del(edge Edge) {
	if edge == nil || edge.Id() == nil { return }
	id := string(edge.Id()[:])
	set.Lock()
	defer set.Unlock()
	delete(set.atommap, id)
	return
}

func (set *EdgeSet) Contains(edge Edge) bool {
	if edge == nil || edge.Id() == nil { return false}
	if _, ok := set.atommap[string(edge.Id()[:])]; ok {
		return true
	}
	return false
}

func (set *EdgeSet) Members() []Edge {
	atoms := []Edge{}
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

func (set *EdgeSet) Count() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.atommap)
}

func (set *EdgeSet) Equal(other *EdgeSet) bool {

	if set.Count() != other.Count() {
		return false
	}
	set.RLock()
	defer set.RUnlock()

	for _, elem := range set.atommap {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}
