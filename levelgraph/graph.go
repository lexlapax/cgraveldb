package levelgraph


import (
		"bytes"
		"fmt"
		"errors"
		"os"
		"encoding/binary"
		"github.com/jmhodges/levigo"
		"path"
		//"github.com/lexlapax/graveldb/core"
)

// Errors in addition to IO errors.
var (
	NoDirectory     = errors.New("need to pass a valid path for a db directory")
	NilValue     = errors.New("nil value passed in argument")
	KeyExists = errors.New("key exists in database")
	KeyDoesNotExist = errors.New("key does not exist in database")
	metadb = "meta.db"	
	elementdb = "element.db"	
	hsdb = "hs.db"	
	propdb = "prop.db"	
)

type DBGraph struct {
	meta *levigo.DB
	elements *levigo.DB
	hs *levigo.DB
	props *levigo.DB
	dbdir string
	opts *levigo.Options
	ro *levigo.ReadOptions
	wo *levigo.WriteOptions
}

func opengraph(dbdir string) (*DBGraph, error) {
	if dbdir == "" {
		return nil, NoDirectory
	}
	err := os.MkdirAll(dbdir, 0755)
	if err != nil {
		return nil, err
	}
	db := new(DBGraph)
	db.dbdir = dbdir
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3<<30))
	opts.SetCreateIfMissing(true)
	filter := levigo.NewBloomFilter(10)
	opts.SetFilterPolicy(filter)

	db.ro = levigo.NewReadOptions()
	db.wo = levigo.NewWriteOptions()

	meta, err := levigo.Open(path.Join(dbdir, metadb), opts)
	if err != nil {return nil, err}
	db.meta = meta
	elements, err := levigo.Open(path.Join(dbdir, elementdb), opts)
	if err != nil {return nil, err}
	db.elements = elements
	hs, err := levigo.Open(path.Join(dbdir, hsdb), opts)
	if err != nil {return nil, err}
	db.hs = hs
	props, err := levigo.Open(path.Join(dbdir, propdb), opts)
	if err != nil {return nil, err}
	db.props = props
	db.keepcount(VertexType, 0)
	db.keepcount(EdgeType, 0)

	return db, err
}

func OpenGraph(dbdir string) (*DBGraph, error ) {
	return opengraph(dbdir)
}

func (db *DBGraph) Close() (bool, error) {
	db.meta.Close()
	db.elements.Close()
	db.hs.Close()
	db.props.Close()
	return true, nil
}

func (db *DBGraph) String() (string) {
	str := fmt.Sprintf("<DBGraph:dbdir=%v>",db.dbdir)
	return str
}


func (db *DBGraph) AddVertex(id []byte) (*DBVertex, error) {
	var vertex *DBVertex = new(DBVertex)
	//vertex := new(DBVertex)
	if id == nil {return nil, NilValue}
	val,err := db.elements.Get(db.ro, id)
	if val != nil {
		return nil, KeyExists
	}
	err = db.elements.Put(db.wo, id, []byte(VertexType))
	if err != nil {return nil, err}
	vertex.id = id
	vertex.Elementtype = VertexType
	vertex.Db = db
	db.keepcount(VertexType, 1)
	return vertex, nil
}

func (db *DBGraph) Vertex(id []byte) *DBVertex {
	if id == nil {return nil }
	val,err := db.elements.Get(db.ro, id)
	if err != nil {return nil}
	if ElementType(val) != VertexType {return nil}
	vertex := &DBVertex{DBElement{db, id, VertexType}}
	return vertex
}

func (db *DBGraph) DelVertex(vertex *DBVertex) error {
	if vertex == nil {	return NilValue }
	id := vertex.Id()
	if id == nil {	return NilValue }
	val,err := db.elements.Get(db.ro, id)
	if err != nil {return err}
	if val == nil {return KeyDoesNotExist}
	err = db.elements.Delete(db.wo, id)
	if err != nil {return err}
	
	db.keepcount(VertexType, -1)
	
	// delete all properties data 
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.props.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	prefix := append(id, []byte("::")...)
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

	// todo - delete all hexastore data

	return nil
}

func (db *DBGraph) Vertices() []*DBVertex {
	vertii := []*DBVertex{}
	var vertex *DBVertex
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.elements.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	it.SeekToFirst()
	for it = it; it.Valid(); it.Next() {
		if ElementType(it.Value()) == VertexType {
			vertex = new(DBVertex)
			vertex.id = it.Key()
			vertex.Elementtype = VertexType
			vertex.Db = db
			vertii = append(vertii, vertex)
		}
	}
	return vertii
}

func (db *DBGraph) keepcount(etype ElementType, upordown int) (uint64) {
	var storedcount, returncount uint64
	var key []byte
	if etype == VertexType {
		key = []byte("vertii")
	} else {
		key = []byte("edges")
	}
	wb := levigo.NewWriteBatch()
	defer wb.Close()
	val, _ := db.meta.Get(db.ro,key)
	if val == nil {
		storedcount = 0
		upordown = 0
	} else {
		storedcount, _ = binary.Uvarint(val)
	}
	switch upordown {
		case -1:
			returncount = storedcount - 1
		case 1:
			returncount = storedcount + 1
		default:
			returncount = storedcount
	}
	if returncount != storedcount || val == nil {
		bufsize := binary.Size(returncount)
		buf := make([]byte, bufsize)
		binary.PutUvarint(buf, returncount)
		_ = db.meta.Put(db.wo, key, buf)
	}
	return returncount
}

func (db *DBGraph) AddEdge(id []byte, outvertex *DBVertex, invertex *DBVertex, label string) (*DBEdge, error) {
	return nil, nil
}

func (db *DBGraph) Edge(id []byte) *DBEdge {
	return nil
}
func (db *DBGraph) DelEdge(edge *DBEdge) error {
	return nil
}
func (db *DBGraph) Edges() []*DBEdge {
	return nil
}

func (db *DBGraph) EdgeCount() uint64 {
	return db.keepcount(EdgeType, 0)
}

func (db *DBGraph) VertexCount() uint64 {
	return db.keepcount(VertexType, 0)
}

