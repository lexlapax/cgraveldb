package levelgraph


import (
		//"bytes"
		"fmt"
		"errors"
		"os"
		"github.com/jmhodges/levigo"
		"github.com/lexlapax/graveldb/core"
)

// Errors in addition to IO errors.
var (
	NoDirectory     = errors.New("need to pass a valid path for a db directory")
	NilValue     = errors.New("nil value passed in argument")
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

func OpenGraph(dbdir string) (core.Graph, error ) {
	return opengraph(dbdir)
}

func (db *DBGraph) String() (string) {
	str := fmt.Sprintf("dbdir=%v",db.dbdir)
	return str
}


func (db *DBGraph) AddVertex(id []byte) (core.Vertex, error) {
	if id == nil {return nil, NilValue}
	
	return nil, nil
}

func (db *DBGraph) Vertex(id []byte) core.Vertex {
	return nil	
}
func (db *DBGraph) DelVertex(vertex core.Vertex) error {
	return nil
}
func (db *DBGraph) Vertices() []core.Vertex {
	return nil	
}

func (db *DBGraph) AddEdge(id []byte, outvertex core.Vertex, invertex core.Vertex, label string) (core.Edge, error) {
	return nil, nil	
}

func (db *DBGraph) Edge(id []byte) core.Edge {
	return nil	
}
func (db *DBGraph) DelEdge(edge core.Edge) error {
	return nil	
}
func (db *DBGraph) Edges() []core.Edge {
	return nil
}
func (db *DBGraph) EdgeCount() uint {
	return 0
}
func (db *DBGraph) VertexCount() uint {
	return 0
}
func (db *DBGraph) Close() (bool, error) {
	return false, nil
}