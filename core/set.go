package core

import (
		"sync"
		"sort"
		"fmt"
)

type StringSet struct {
	smap  map[string]int
	sync.RWMutex
}

func NewStringSet() *StringSet {
	set := new(StringSet)
	set.smap = make(map[string]int)
	return set
}

func (set *StringSet) Add(s string) {
	if s == "" { return }
	set.Lock()
	defer set.Unlock()
	if _, ok := set.smap[s]; ok {
		return
	} else {
		set.smap[s] = 1
	}
	return
}

func (set *StringSet) Del(s string) {
	if s == "" { return }
	set.Lock()
	defer set.Unlock()
	delete(set.smap, s)
	return
}

func (set *StringSet) Contains(s string) bool {
	if s == "" { return false }
	if _, ok := set.smap[s]; ok {
		return true
	}
	return false
}

func (set *StringSet) Members() []string {
	keys := []string{}
	set.RLock()
	defer set.RUnlock()
	for k, _ := range set.smap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (set *StringSet) Count() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.smap)
}

func (set *StringSet) Equal(other *StringSet) bool {

	if set.Count() != other.Count() {
		return false
	}
	set.RLock()
	defer set.RUnlock()

	for k, _ := range set.smap {
		if !other.Contains(k) {
			return false
		}
	}
	return true
}

func (set *StringSet) String() string {
	return "&{Stringset[] " + fmt.Sprintf("%v", set.Members()) + "}"
}

//-----atomSet
type AtomSet struct {
	atommap map[string]Atom
	sync.RWMutex
}

func NewAtomSet() *AtomSet {
	set := new(AtomSet)
	set.atommap = make(map[string]Atom)
	return set
}

func (set *AtomSet) Add(atom Atom) {
	if atom == nil || atom.Id() == nil { return }
	id := string(atom.Id()[:])
	set.Lock()
	defer set.Unlock()
	if _, ok := set.atommap[id]; ok {
		return
	} else {
		set.atommap[id] = atom
	}
	return
}

func (set *AtomSet) Del(atom Atom) {
	if atom == nil || atom.Id() == nil { return }
	id := string(atom.Id()[:])
	set.Lock()
	defer set.Unlock()
	delete(set.atommap, id)
	return
}

func (set *AtomSet) Contains(atom Atom) bool {
	if atom == nil || atom.Id() == nil { return false}
	if _, ok := set.atommap[string(atom.Id()[:])]; ok {
		return true
	}
	return false
}

func (set *AtomSet) Members() []Atom {
	atoms := []Atom{}
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

func (set *AtomSet) Count() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.atommap)
}

func (set *AtomSet) Clear() {
	set.Lock()
	defer set.Unlock()
	set.atommap = make(map[string]Atom)
}

func (set *AtomSet) Equal(other *AtomSet) bool {

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

// //-----edgeSet
// type EdgeSet struct {
// 	edgemap map[string]Edge
// 	sync.RWMutex
// }

// func NewEdgeSet() *EdgeSet {
// 	set := new(EdgeSet)
// 	set.edgemap = make(map[string]Edge)
// 	return set
// }

// func (set *EdgeSet) Add(edge Edge) {
// 	if edge == nil || edge.Id() == nil { return }
// 	id := string(edge.Id()[:])
// 	set.Lock()
// 	defer set.Unlock()
// 	if _, ok := set.edgemap[id]; ok {
// 		return
// 	} else {
// 		set.edgemap[id] = edge
// 	}
// 	return
// }

// func (set *EdgeSet) Del(edge Edge) {
// 	if edge == nil || edge.Id() == nil { return }
// 	id := string(edge.Id()[:])
// 	set.Lock()
// 	defer set.Unlock()
// 	delete(set.edgemap, id)
// 	return
// }

// func (set *EdgeSet) Contains(edge Edge) bool {
// 	if edge == nil || edge.Id() == nil { return false}
// 	if _, ok := set.edgemap[string(edge.Id()[:])]; ok {
// 		return true
// 	}
// 	return false
// }

// func (set *EdgeSet) Members() []Edge {
// 	atoms := []Edge{}
// 	keys := []string{}
// 	set.RLock()
// 	defer set.RUnlock()
// 	for k, _ := range set.edgemap {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)
// 	for _, k := range keys {
// 		atoms = append(atoms, set.edgemap[k])
// 	}
// 	return atoms
// }

// func (set *EdgeSet) Count() int {
// 	set.RLock()
// 	defer set.RUnlock()
// 	return len(set.edgemap)
// }

// func (set *EdgeSet) Equal(other *EdgeSet) bool {

// 	if set.Count() != other.Count() {
// 		return false
// 	}
// 	set.RLock()
// 	defer set.RUnlock()

// 	for _, elem := range set.edgemap {
// 		if !other.Contains(elem) {
// 			return false
// 		}
// 	}
// 	return true
// }
