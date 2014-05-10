package mem

import (
		"github.com/lexlapax/graveldb/core"
		"sync"
		"sort"
)

type AtomMem struct {
	graph *GraphMem
	id string
	atomType core.AtomType
	props map[string][]byte 
	sync.RWMutex
}

func NewAtomMem(graph *GraphMem, id string, atomType core.AtomType) *AtomMem {
	atom := new(AtomMem)
	atom.id = id 
	atom.graph = graph
	atom.atomType = atomType
	atom.props = make(map[string][]byte)
	return atom
}

func (atom *AtomMem) Id() string {
	return atom.id
}

func (atom *AtomMem) Property(prop string) ([]byte, error) {
	atom.RLock()
	defer atom.RUnlock()
 	if val, ok := atom.props[prop]; ok {
        return val, nil
    }
	return nil, nil
}

func (atom *AtomMem) SetProperty(prop string, value []byte) error {
	atom.Lock()
	defer atom.Unlock()
	oldvalue := atom.props[prop]
	atom.props[prop] = value
	atom.graph.keyindex.update(prop, string(value[:]), string(oldvalue[:]), string(atom.id[:]), atom.atomType)
	return nil
}

func (atom *AtomMem) DelProperty(prop string) error {
	atom.Lock()
	defer atom.Unlock()
	oldvalue := atom.props[prop]
 	delete(atom.props, prop)
 	atom.graph.keyindex.remove(prop, string(oldvalue[:]), string(atom.id[:]), atom.atomType)
	return nil
}

func (atom *AtomMem) PropertyKeys() ([]string, error) {
	atom.RLock()
	defer atom.RUnlock()
	keys := []string{}
	for k := range atom.props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}
