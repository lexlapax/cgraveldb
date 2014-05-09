package levigo

import (
		"os"
		"path"
		"strings"
		"github.com/lexlapax/graveldb/core"
		levigoindex "github.com/lexlapax/graveldb/util/levigo"
		// "sync"
		//"strconv"
		// "fmt"
)

func NewKeyIndex(dbdir string) *KeyIndex {
	index := &KeyIndex{}
	index.Open(dbdir)
	return index
}
var ( 
	indexdir = "indices"
)

type KeyIndex struct {
	indexc *levigoindex.IndexContainer
	dbdir string
	isopen bool
}

func (index *KeyIndex) Open(args ...interface{}) error {
	if index.isopen == true { return nil }
	if len(args) > 0 {
		if aString, found := args[0].(string); found {
			if aString == "" {
				return NoDirectory
			} else {
				index.dbdir = path.Join(aString, indexdir)
			}
		} else {
			return InvalidParameterValue
		}
	}
	if index.dbdir == "" { return NoDirectory }
	err := os.MkdirAll(index.dbdir, 0755)
	if err != nil { return err }
	indexc, err :=  levigoindex.OpenIndexContainer(index.dbdir)
	if err != nil { return err }
	index.indexc = indexc
	index.isopen = true
	return nil
}

func (index *KeyIndex) Close() error {
	index.isopen = false
	index.indexc.Close()
	return nil
}

func (index *KeyIndex) Clear() error {
	index.indexc.Clear()
	return nil
}

func (index *KeyIndex) update(propkey string, newValue string, oldValue string, atomid string, atomType core.AtomType) error {
	if propkey == "" || atomid == "" { return core.ErrNilValue }
	if newValue == "" { return nil}
	var key string
	if atomType == core.VertexType {
		key = "v" + fieldsep + propkey
	} else if atomType == core.EdgeType {
		key = "e" + fieldsep + propkey
	} else {
		return nil
	}
	if index.indexc.HasIndex(key) {
		index.indexc.DelDoc(key, atomid)
		index.indexc.AddDoc(key, atomid, newValue)
	}
	return nil
}

func (index *KeyIndex) remove(propkey string, oldValue string, atomid string, atomType core.AtomType) error {
	if propkey == "" || atomid == "" { return core.ErrNilValue }
	var key string
	if atomType == core.VertexType {
		key = "v" + fieldsep + propkey
	} else if atomType == core.EdgeType {
		key = "e" + fieldsep + propkey
	} else {
		return nil
	}
	if index.indexc.HasIndex(key) {
		index.indexc.DelDoc(key, atomid)
	}
	return nil
}


// type KeyIndex interface
func (index *KeyIndex) createKeyIndex(key string, atomType core.AtomType) error {
	if key == "" { return core.ErrNilValue }
	var indexkey string
	if atomType == core.VertexType {
		indexkey = "v" + fieldsep + key
	} else if atomType == core.EdgeType {
		indexkey = "e" + fieldsep + key
	} else {
		return nil
	}
	//fmt.Printf("v=%v\n", indexkey)
	_, err := index.indexc.AddIndex(indexkey)
	return err
}

func (index *KeyIndex) dropKeyIndex(key string, atomType core.AtomType) error {
	if key == "" { return core.ErrNilValue }
	var indexkey string
	if atomType == core.VertexType {
		indexkey = "v" + fieldsep + key
	} else if atomType == core.EdgeType {
		indexkey = "e" + fieldsep + key
	} else {
		return nil
	}
	err := index.indexc.DelIndex(indexkey)
	return err
}

func (index *KeyIndex) indexedKeys(atomType core.AtomType) []string {
	vkeys := []string{}
	ekeys := []string{}
	indices, _ := index.indexc.Indices()
	// fmt.Printf("v=%v\n", indices)
	if len(indices) < 1 { return []string{} }
	for _, k := range indices {
		args := strings.Split(k, fieldsep)
		if len(args) != 2 { continue }
		if args[0] == "v" { 
			vkeys = append(vkeys, args[1])
		} else if args[0] == "e" {
			ekeys = append(ekeys, args[1])
		}
	}

	if atomType == core.VertexType {
		return vkeys
	} else if atomType == core.EdgeType {
		return ekeys
	}
	return []string{}
}

func (index *KeyIndex) searchIds(key string, value string, atomType core.AtomType) []string {
	ids := []string{}
	searchkey := ""
	if key == "" || value == "" { return ids }
	if atomType == core.VertexType {
		searchkey = "v" + fieldsep + key
	} else if atomType == core.EdgeType {
		searchkey = "e" + fieldsep + key
	} else {
		return ids
	}
	ids = index.indexc.Search(searchkey, value)
	return ids
}
