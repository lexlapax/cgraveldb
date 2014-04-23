package mem

import (
		//"github.com/lexlapax/graveldb/core"
		"sync"
)

type NodeType string
const (
	VertexType NodeType = "1"
	EdgeType ="2"
)

type AtomMem struct {
	graph *GraphMem
	id []byte
	nodeType NodeType
	props map[string][]byte 
	sync.RWMutex
}

func NewAtomMem(graph *GraphMem, id []byte, nodeType NodeType) *AtomMem {
	atom := new(AtomMem)
	atom.id = id 
	atom.graph = graph
	atom.nodeType = nodeType
	atom.props = make(map[string][]byte)
	return atom
}

func (atom *AtomMem) Id() []byte {
	return atom.Id()
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
	atom.props[prop] = value
	return nil
}

func (atom *AtomMem) DelProperty(prop string) error {
	atom.Lock()
	defer atom.Unlock()
 	delete(atom.props, prop)
	return nil
}

func (atom *AtomMem) PropertyKeys() ([]string, error) {
	atom.RLock()
	defer atom.RUnlock()
	keys := []string{}
	for k := range atom.props {
		keys = append(keys, k)
	}
	return keys, nil
}
