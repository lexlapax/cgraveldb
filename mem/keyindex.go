package mem

import (
		"github.com/lexlapax/graveldb/core"
		"sync"
		//"strconv"
		//"fmt"
)

type GraphKeyIndex struct {
	vertexkeys map[string]int 
	edgekeys map[string]int
	vlock *sync.RWMutex 
	elock *sync.RWMutex
}

// type KeyIndexableGraph interface
func (graph *GraphMem) CreateKeyIndex(key string, atomType core.AtomType) error {
	if key == "" { return core.ErrNilValue }
	return nil
}

func (graph *GraphMem) DropKeyIndex(key string, atomType core.AtomType) error {
	if key == "" { return core.ErrNilValue }
	return nil
}

func (graph *GraphMem) IndexedKeys(atomType core.AtomType) []string {
	return nil
}

