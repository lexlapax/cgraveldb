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
	edgedb = "edge.db"	
	hexaindexdb = "hexaindex.db"
	propdb = "prop.db"
	recsep = "\x1f" //the ascii unit separator dec val 31
	propVertiiCount = "vertiicount"
	propEdgeCount = "edgecount"
	propRecSep = "recsep"
)

type DBGraph struct {
	meta *levigo.DB
	elements *levigo.DB
	edges *levigo.DB
	hexaindex *levigo.DB
	props *levigo.DB
	dbdir string
	opts *levigo.Options
	ro *levigo.ReadOptions
	wo *levigo.WriteOptions
	recsep []byte
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
	db.recsep = []byte(recsep)

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
	recsepbytes, _ := db.getDbProperty(propRecSep)
	if recsepbytes == nil {
		recsepbytes = []byte(recsep)
		_,err = db.putDbProperty(propRecSep,recsepbytes)
		if err != nil {return nil, err}
	}
	db.recsep = recsepbytes

	elements, err := levigo.Open(path.Join(dbdir, elementdb), opts)
	if err != nil {return nil, err}
	db.elements = elements

	edges, err :=  levigo.Open(path.Join(dbdir, edgedb), opts)
	if err != nil {return nil, err}
	db.edges = edges
	
	hexaindex, err := levigo.Open(path.Join(dbdir, hexaindexdb), opts)
	if err != nil {return nil, err}
	db.hexaindex = hexaindex
	
	props, err := levigo.Open(path.Join(dbdir, propdb), opts)
	if err != nil {return nil, err}
	db.props = props
	
	db.keepcount(VertexType, 0)
	db.keepcount(EdgeType, 0)

	return db, nil
}

func OpenGraph(dbdir string) (*DBGraph, error ) {
	return opengraph(dbdir)
}

func (db *DBGraph) Close() (bool, error) {
	db.meta.Close()
	db.elements.Close()
	db.edges.Close()
	db.hexaindex.Close()
	db.props.Close()
	return true, nil
}

func (db *DBGraph) String() (string) {
	str := fmt.Sprintf("#DBGraph:dbdir=%v#",db.dbdir)
	return str
}

func(db *DBGraph) getDbProperty(prop string) ([]byte, error){
	if prop == "" {return nil, NilValue}
	val, err := db.meta.Get(db.ro, []byte(prop))
	if err != nil {return nil, err}
	return val, nil
}

func(db *DBGraph) putDbProperty(prop string, val []byte) ([]byte, error){
	if prop == "" {return nil, NilValue}
	key := []byte(prop)
	oldval, err := db.meta.Get(db.ro, key)
	if err != nil {return nil, err}
	err2 := db.meta.Put(db.wo, key, val)
	if err2 != nil {return nil, err}
	return oldval, nil
}


func (db *DBGraph) keepcount(etype ElementType, upordown int) (uint64) {
	var storedcount, returncount uint64
	var key string
	if etype == VertexType {
		key = propVertiiCount
	} else {
		key = propEdgeCount
	}
	//wb := levigo.NewWriteBatch()
	//defer wb.Close()
	val, _ := db.getDbProperty(key)
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
		db.putDbProperty(key, buf)
	}
	return returncount
}

func (db *DBGraph) AddVertex(id []byte) (*DBVertex, error) {
	if id == nil {return nil, NilValue}
	val, err := db.elements.Get(db.ro, id)
	if val != nil {
		return nil, KeyExists
	}
	err = db.elements.Put(db.wo, id, []byte(VertexType))
	if err != nil {return nil, err}
	vertex := &DBVertex{&DBElement{db,id,VertexType}}
	db.keepcount(VertexType, 1)
	return vertex, nil
}

func (db *DBGraph) Vertex(id []byte) *DBVertex {
	if id == nil {return nil }
	val,err := db.elements.Get(db.ro, id)
	if err != nil {return nil}
	if val == nil {return nil}
	if ElementType(val) != VertexType {return nil}
	vertex := &DBVertex{&DBElement{db, id, VertexType}}
	/*
	vertex := new(DBVertex)
	vertex.db = db
	vertex.id = id
	vertex.Elementtype = VertexType
	*/
	return vertex
}

func (db *DBGraph) DelVertex(vertex *DBVertex) error {
	if vertex == nil {	return NilValue }
	if vertex.DBElement == nil { return NilValue }
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
			vertex = &DBVertex{&DBElement{db, it.Key(), VertexType}}
			vertii = append(vertii, vertex)
		}
	}
	return vertii
}

func (db *DBGraph) AddEdge(id []byte, outvertex *DBVertex, invertex *DBVertex, label string) (*DBEdge, error) {
	if (id == nil) {return nil, NilValue}
	if (outvertex == nil) {return nil, NilValue}
	if (invertex == nil) {return nil, NilValue}

	val,err := db.elements.Get(db.ro, id)
	if val != nil {
		return nil, KeyExists
	}
	err = db.elements.Put(db.wo, id, []byte(EdgeType))
	if err != nil {return nil, err}
	edge := &DBEdge{&DBElement{db, id, EdgeType}, outvertex, invertex, label}

	//fmt.Printf("evin=%v\n", edgerecord)
	err = db.edges.Put(db.wo, id, db.toEdgeRecord(outvertex, invertex, label))
	if err != nil {return nil, err}
	//todo - add hexascale index
	db.keepcount(EdgeType, 1)
	return edge, nil
}

func (db *DBGraph) toEdgeRecord(outvertex *DBVertex, invertex *DBVertex, label string) ([]byte) {
	edgevalues := [][]byte{}
	edgevalues = append(edgevalues,outvertex.id, invertex.id, []byte(label))
	edgerecord := bytes.Join(edgevalues, db.recsep)
	return edgerecord
}

func (db *DBGraph) fromEdgeRecord(record []byte) (*DBVertex, *DBVertex, string) {
	if record == nil { return nil, nil, ""}
	edgevalues := bytes.Split(record, db.recsep)

	outvertex := db.Vertex(edgevalues[0])
	invertex := db.Vertex(edgevalues[1])
	label := string(edgevalues[2][:])
	return outvertex, invertex, label
}

func (db *DBGraph) Edge(id []byte) *DBEdge {
	if id == nil {return nil }
	val,err := db.elements.Get(db.ro, id)
	if err != nil {return nil}
	if val == nil {return nil}
	if ElementType(val) != EdgeType {return nil}
	val, err = db.edges.Get(db.ro, id)
	if err != nil {return nil}
	if val == nil {return nil}
	outvertex, invertex, label := db.fromEdgeRecord(val)
	//fmt.Printf("evout=%v\n", val)
	edge := &DBEdge{&DBElement{db, id, EdgeType}, outvertex, invertex, label}
	return edge
}

func (db *DBGraph) DelEdge(edge *DBEdge) error {
	if edge == nil {	return NilValue }
	if edge.DBElement == nil { return NilValue }
	id := edge.Id()
	if id == nil {	return NilValue }
	val,err := db.elements.Get(db.ro, id)
	if err != nil {return err}
	if val == nil {return KeyDoesNotExist}
	err = db.elements.Delete(db.wo, id)
	if err != nil {return err}
	
	err = db.edges.Delete(db.wo, id)
	db.keepcount(EdgeType, -1)
	
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

	// todo - delete all hexastore data

	return nil
}

func (db *DBGraph) Edges() []*DBEdge {
	edges := []*DBEdge{}
	var edge *DBEdge
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.edges.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	it.SeekToFirst()
	var label string
	var outvertex, invertex *DBVertex
	for it = it; it.Valid(); it.Next() {
		outvertex, invertex, label = db.fromEdgeRecord(it.Value())

		edge = &DBEdge{&DBElement{db, it.Key(), EdgeType}, outvertex, invertex, label}
		edges = append(edges, edge)
	}
	return edges
}

func (db *DBGraph) EdgeCount() uint64 {
	return db.keepcount(EdgeType, 0)
}

func (db *DBGraph) VertexCount() uint64 {
	return db.keepcount(VertexType, 0)
}

