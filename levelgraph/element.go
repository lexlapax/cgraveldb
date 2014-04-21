package levelgraph


import (
		"bytes"
		//"fmt"
		//"errors"
		//"os"
		"github.com/jmhodges/levigo"
		//"github.com/lexlapax/graveldb/core"
)

type ElementType string
const (
	VertexType ElementType = "1"
	EdgeType ="2"
)

type DBElement struct {
	db *DBGraph
	id []byte
	Elementtype ElementType
}

func (db *DBGraph) getPropKey(id []byte, prop string) ([]byte) {
	keyvalues := [][]byte{}
	keyvalues = append(keyvalues,id, []byte(prop))
	key := bytes.Join(keyvalues, db.recsep)
	return key
}

func (element *DBElement) Property(prop string) ([]byte) {
	if prop == "" { return nil }
	key := element.db.getPropKey(element.id, prop)
	val, err := element.db.props.Get(element.db.ro, key)
	if err != nil {return nil}
	return val
}

func (element *DBElement) SetProperty(prop string, value []byte) (error){
	if prop == "" { return nil }
	key := element.db.getPropKey(element.id, prop)
	err := element.db.props.Put(element.db.wo, key, value)
	return err

}

func (element *DBElement) DelProperty(prop string) ([]byte) {
	if prop == "" { return nil}
	key := element.db.getPropKey(element.id, prop)
	val, err := element.db.props.Get(element.db.ro, key)
	if err != nil {return nil}
	err = element.db.props.Delete(element.db.wo, key)
	if err != nil {return nil}
	return val
}

func (element *DBElement) PropertyKeys() ([]string) {
	propkeys := []string{}
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := element.db.props.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	prefix := append(element.id, element.db.recsep...)
	it.Seek(prefix)
	var prop []byte
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		prop = bytes.Split(it.Key(), element.db.recsep)[1]
		propkeys = append(propkeys, string(prop[:]))
	}
	return propkeys
}

func (element *DBElement) Id() ([]byte) {
	return element.id
}

func (element *DBElement) IdAsString() (string) {
	return string(element.id[:])
}

/*

	// delete all properties data 
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.props.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	prefix := append(id, db.recsep...)
	it.Seek(prefix)
	propkeys := [][]byte{}
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		propkeys = append(propkeys, it.Key())
	}
	wb := levigo.NewWriteBatch()
	defer wb.Close()
	for _, propkey := range propkeys {
		wb.Delete(propkey)
	}
	err = db.props.Write(db.wo, wb)
	if err != nil {return err}
*/