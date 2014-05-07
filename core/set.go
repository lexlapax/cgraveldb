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

func (set *StringSet) AddArray(sarray []string) {
	if sarray == nil { return }
	set.Lock()
	defer set.Unlock()
	for _, s := range sarray {
		if s == "" { continue }
		if _, ok := set.smap[s]; ok {
			continue
		} else {
			set.smap[s] = 1
		}
	}
	return
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

func (set *StringSet) Clear() {
	set.Lock()
	defer set.Unlock()
	set.smap = make(map[string]int)
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


//-----bytearraySet
type ByteArraySet struct {
	arraymap map[string][]byte
	sync.RWMutex
}

func NewByteArraySet() *ByteArraySet {
	set := new(ByteArraySet)
	set.arraymap = make(map[string][]byte)
	return set
}

func (set *ByteArraySet) Add(ba []byte) {
	if ba == nil { return }
	id := string(ba[:])
	set.Lock()
	defer set.Unlock()
	if _, ok := set.arraymap[id]; ok {
		return
	} else {
		set.arraymap[id] = ba
	}
	return
}

func (set *ByteArraySet) Del(ba []byte) {
	if ba == nil { return }
	id := string(ba[:])
	set.Lock()
	defer set.Unlock()
	delete(set.arraymap, id)
	return
}

func (set *ByteArraySet) Contains(ba []byte) bool {
	if ba == nil { return false}
	if _, ok := set.arraymap[string(ba[:])]; ok {
		return true
	}
	return false
}

func (set *ByteArraySet) Members() [][]byte {
	members := [][]byte{}
	keys := []string{}
	set.RLock()
	defer set.RUnlock()
	for k, _ := range set.arraymap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		members = append(members, set.arraymap[k])
	}
	return members
}

func (set *ByteArraySet) Count() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.arraymap)
}

func (set *ByteArraySet) Clear() {
	set.Lock()
	defer set.Unlock()
	set.arraymap = make(map[string][]byte)
}

func (set *ByteArraySet) Equal(other *ByteArraySet) bool {

	if set.Count() != other.Count() {
		return false
	}
	set.RLock()
	defer set.RUnlock()

	for _, elem := range set.arraymap {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}
