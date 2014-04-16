package levelgraph


import (
		//"bytes"
		"fmt"
		"errors"
		"os"
		"github.com/jmhodges/levigo"
		//"github.com/lexlapax/graveldb/core"
)

// Errors in addition to IO errors.
var (
	NoDirectory     = errors.New("need to pass a valid path for a db directory")
	NilValue     = errors.New("nil value passed in argument")
	KeyExists = errors.New("key exists in database")
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

	meta, err := levigo.Open(dbdir + "/" + metadb, opts)
	if err != nil {return nil, err}
	db.meta = meta
	elements, err := levigo.Open(dbdir + "/" + elementdb, opts)
	if err != nil {return nil, err}
	db.elements = elements
	hs, err := levigo.Open(dbdir + "/" + hsdb, opts)
	if err != nil {return nil, err}
	db.hs = hs
	props, err := levigo.Open(dbdir + "/" + propdb, opts)
	if err != nil {return nil, err}
	db.props = props

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
	str := fmt.Sprintf("dbdir=%v",db.dbdir)
	return str
}


func (db *DBGraph) AddVertex(id []byte) (*DBVertex, error) {
	var vertex *DBVertex = new(DBVertex)
	//vertex := new(DBVertex)
	vertex.id = id
	vertex.elementtype = VertexType
	vertex.db = db
	if id == nil {return vertex, NilValue}
	val,err := db.elements.Get(db.ro, id)
	if val != nil {
		return vertex, KeyExists
	}
	err = db.elements.Put(db.wo, id, []byte(VertexType))
	if err != nil {return vertex, err}
	return vertex, nil
}

func (db *DBGraph) Vertex(id []byte) *DBVertex {
	return nil //DBVertex{}
}
func (db *DBGraph) DelVertex(vertex *DBVertex) error {
	return nil
}
func (db *DBGraph) Vertices() []*DBVertex {
	return nil
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

func (db *DBGraph) EdgeCount() uint {
	return 0
}
func (db *DBGraph) VertexCount() uint {
	return 0
}

