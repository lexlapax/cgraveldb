package mem

import (
		"github.com/lexlapax/graveldb/core"
		"sync"
		//"strconv"
		//"fmt"
)

var (
	registerindex sync.Once
)

func init() {
	registerindex.Do(func() {core.Register("memindexed", &KeyIndexedGraph{})} )
}

func NewKeyIndexedGraph() core.KeyIndexableGraph {
	graph := &KeyIndexedGraph{}
	graph.Open()
	return graph
}

type KeyIndexedGraph struct {
	GraphMem
	vIdx map[string]*InvertedIndex 
	eIdx map[string]*InvertedIndex
	vidxLk *sync.RWMutex 
	eidxLk *sync.RWMutex
}


// type KeyIndexableGraph interface
func (graph *KeyIndexedGraph) CreateKeyIndex(key string, atomType core.AtomType) error {
	if key == "" { return core.ErrNilValue }
	if atomType == core.VertexType {
		graph.vidxLk.Lock()
		defer graph.vidxLk.Unlock()
		if _, ok := graph.vIdx[key]; ok { return nil }
		graph.vIdx[key] = NewInvertedIndex()
	} else if atomType == core.EdgeType {
		graph.eidxLk.Lock()
		defer graph.eidxLk.Unlock()
		if _, ok := graph.eIdx[key]; ok { return nil }
		graph.eIdx[key] = NewInvertedIndex()
	}
	return nil
}

func (graph *KeyIndexedGraph) DropKeyIndex(key string, atomType core.AtomType) error {
	if key == "" { return core.ErrNilValue }
	if atomType == core.VertexType {
		graph.vidxLk.Lock()
		defer graph.vidxLk.Unlock()
		delete(graph.vIdx, key)
	} else if atomType == core.EdgeType {
		graph.eidxLk.Lock()
		defer graph.eidxLk.Unlock()
		delete(graph.eIdx, key)
	}
	return nil
}

func (graph *KeyIndexedGraph) IndexedKeys(atomType core.AtomType) []string {
	keys := []string{}
	if atomType == core.VertexType {
		graph.vidxLk.RLock()
		defer graph.vidxLk.RUnlock()
		for k, _ := range graph.vIdx {
			keys = append(keys, k)
		}
	} else if atomType == core.EdgeType {
		graph.eidxLk.RLock()
		defer graph.eidxLk.RUnlock()
		for k, _ := range graph.eIdx {
			keys = append(keys, k)
		}
	}
	return keys
}

func (graph *KeyIndexedGraph) Open(args ...interface{}) error {
	graph.vIdx = make(map[string]*InvertedIndex)
	graph.eIdx = make(map[string]*InvertedIndex)
	graph.vidxLk = &sync.RWMutex{}
	graph.eidxLk = &sync.RWMutex{}
	return graph.GraphMem.Open(args...)
}

func (graph *KeyIndexedGraph) Close() error {
	graph.vIdx = nil
	graph.eIdx = nil
	graph.vidxLk = &sync.RWMutex{}
	graph.eidxLk = &sync.RWMutex{}
	return graph.GraphMem.Close()
}

func (graph *KeyIndexedGraph) Clear() error {
	graph.vidxLk.Lock()
	graph.vIdx = make(map[string]*InvertedIndex)
	graph.vidxLk.Unlock()
	graph.eidxLk.Lock()
	graph.eIdx = make(map[string]*InvertedIndex)
	graph.eidxLk.Unlock()
	return graph.GraphMem.Clear()
}

