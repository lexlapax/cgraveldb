package mem

import (
		"github.com/lexlapax/graveldb/core"
		"sync"
		//"strconv"
		// "fmt"
)

func NewKeyIndex() *KeyIndex {
	index := &KeyIndex{}
	index.open()
	return index
}

type KeyIndex struct {
	vIdx map[string]*InvertedIndex 
	eIdx map[string]*InvertedIndex
	vidxLk *sync.RWMutex 
	eidxLk *sync.RWMutex
	isopen bool
}

func (index *KeyIndex) update(propkey string, newValue string, oldValue string, atomid string, atomType core.AtomType) error {
	if propkey == "" || atomid == "" { return core.ErrNilValue }
	if newValue == "" { return nil}
	if atomType == core.VertexType {
		index.vidxLk.Lock()
		defer index.vidxLk.Unlock()
		if iindex, ok := index.vIdx[propkey]; ok {
			// fmt.Printf("v=%v\n", iindex)
			iindex.DelDoc(atomid)
			iindex.AddDoc(atomid, newValue)
			index.vIdx[propkey] = iindex
		}
	} else if atomType == core.EdgeType {
		// fmt.Printf("k=%v,nv=%v,ov=%v,id=%v\n", propkey, newValue, oldValue, atomid)
		index.eidxLk.Lock()
		defer index.eidxLk.Unlock()
		if iindex, ok := index.eIdx[propkey]; ok {
			iindex.DelDoc(atomid)
			iindex.AddDoc(atomid, newValue)
			index.eIdx[propkey] = iindex
		}
	}
	return nil
}

func (index *KeyIndex) remove(propkey string, oldValue string, atomid string, atomType core.AtomType) error {
	if propkey == "" || atomid == "" { return core.ErrNilValue }
	if atomType == core.VertexType {
		index.vidxLk.Lock()
		defer index.vidxLk.Unlock()
		if iindex, ok := index.vIdx[propkey]; ok {
			iindex.DelDoc(atomid)
			index.vIdx[propkey] = iindex
		}
	} else if atomType == core.EdgeType {
		index.eidxLk.Lock()
		defer index.eidxLk.Unlock()
		if iindex, ok := index.eIdx[propkey]; ok {
			iindex.DelDoc(atomid)
			index.eIdx[propkey] = iindex
		}
	}
	return nil
}


// type KeyIndex interface
func (index *KeyIndex) createKeyIndex(key string, atomType core.AtomType) error {
	if key == "" { return core.ErrNilValue }
	if atomType == core.VertexType {
		index.vidxLk.Lock()
		defer index.vidxLk.Unlock()
		if _, ok := index.vIdx[key]; ok { return nil }
		index.vIdx[key] = NewInvertedIndex()
	} else if atomType == core.EdgeType {
		index.eidxLk.Lock()
		defer index.eidxLk.Unlock()
		if _, ok := index.eIdx[key]; ok { return nil }
		index.eIdx[key] = NewInvertedIndex()
	}
	return nil
}

func (index *KeyIndex) dropKeyIndex(key string, atomType core.AtomType) error {
	if key == "" { return core.ErrNilValue }
	if atomType == core.VertexType {
		index.vidxLk.Lock()
		defer index.vidxLk.Unlock()
		delete(index.vIdx, key)
	} else if atomType == core.EdgeType {
		index.eidxLk.Lock()
		defer index.eidxLk.Unlock()
		delete(index.eIdx, key)
	}
	return nil
}

func (index *KeyIndex) indexedKeys(atomType core.AtomType) []string {
	keys := []string{}
	if atomType == core.VertexType {
		index.vidxLk.RLock()
		defer index.vidxLk.RUnlock()
		for k, _ := range index.vIdx {
			keys = append(keys, k)
		}
	} else if atomType == core.EdgeType {
		index.eidxLk.RLock()
		defer index.eidxLk.RUnlock()
		for k, _ := range index.eIdx {
			keys = append(keys, k)
		}
	}
	return keys
}

func (index *KeyIndex) searchIds(key string, value string, atomType core.AtomType) []string {
	ids := []string{}
	if key == "" || value == "" { return ids }
	if atomType == core.VertexType {
		index.vidxLk.RLock()
		defer index.vidxLk.RUnlock()
		if idx, ok := index.vIdx[key]; ok {
			//fmt.Printf("k=%v,v=%v\n", key, value)
			ids = idx.Search(value)
		} 
	} else if atomType == core.EdgeType {
		index.eidxLk.RLock()
		defer index.eidxLk.RUnlock()
		if idx, ok := index.eIdx[key]; ok {
			ids = idx.Search(value)
		} 
	}
	//fmt.Printf("v=%v\n", ids)
	return ids
}

func (index *KeyIndex) open(args ...interface{}) error {
	if index.isopen == true { return nil }
	index.vIdx = make(map[string]*InvertedIndex)
	index.eIdx = make(map[string]*InvertedIndex)
	index.vidxLk = &sync.RWMutex{}
	index.eidxLk = &sync.RWMutex{}
	index.isopen = true
	return nil
}

func (index *KeyIndex) close() error {
	index.isopen = false
	index.vIdx = nil
	index.eIdx = nil
	index.vidxLk = &sync.RWMutex{}
	index.eidxLk = &sync.RWMutex{}
	return nil
}

func (index *KeyIndex) clear() error {
	index.vidxLk.Lock()
	index.vIdx = make(map[string]*InvertedIndex)
	index.vidxLk.Unlock()
	index.eidxLk.Lock()
	index.eIdx = make(map[string]*InvertedIndex)
	index.eidxLk.Unlock()
	return nil
}

