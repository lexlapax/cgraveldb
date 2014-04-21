package levelgraph


import (
		"bytes"
		"fmt"
		"errors"
		"os"
		"sync"
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
	InvalidParameterValue = errors.New("the value of parameter passed is invalid")
	DBNotOpen = errors.New("db is not open")
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

// * spo::A::C::B
// * sop::A::B::C
// * ops::B::C::A
// * osp::B::A::C
// * pso::C::A::B
// * pos::C::B::A
type HexaIndexType string
const (
	SPO HexaIndexType = "1"
	SOP HexaIndexType = "2"
	OPS HexaIndexType = "3"
	OSP HexaIndexType = "4"
	PSO HexaIndexType = "5"
	POS HexaIndexType = "6"
)

type HexaIndexKeys struct {
	spo []byte
	sop []byte
	ops []byte
	osp []byte
	pso []byte
	pos []byte
}


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
	IsOpen bool
	rwlock *sync.RWMutex
}

func openMeta(dbdir string, ro *levigo.ReadOptions, wo *levigo.WriteOptions, opts *levigo.Options) (*levigo.DB, []byte, error) {
	meta, err := levigo.Open(path.Join(dbdir, metadb), opts)
	if err != nil {return nil, nil, err}
	recsepbytes, _ := meta.Get(ro, []byte(propRecSep))
	if recsepbytes == nil {
		recsepbytes = []byte(recsep)
		err = meta.Put(wo, []byte(propRecSep), recsepbytes) 
		if err != nil {return nil, nil, err}
	}
	return meta, recsepbytes, nil
}

func openElements(dbdir string, opts *levigo.Options) (*levigo.DB, error) {
	return levigo.Open(path.Join(dbdir, elementdb), opts)
}

func openEdges(dbdir string, opts *levigo.Options) (*levigo.DB, error) {
	return levigo.Open(path.Join(dbdir, edgedb), opts)
}

func openHexaIndex(dbdir string, opts *levigo.Options) (*levigo.DB, error) {
	return levigo.Open(path.Join(dbdir, hexaindexdb), opts)
}

func openProps(dbdir string, opts *levigo.Options) (*levigo.DB, error) {
	return levigo.Open(path.Join(dbdir, propdb), opts)
}

func (db *DBGraph) Open() (error) {
	if db.IsOpen == true {return errors.New("db already open") }
	err := os.MkdirAll(db.dbdir, 0755)
	if err != nil { return err }

	//db.recsep = []byte(recsep)
	db.rwlock = &sync.RWMutex{}
	db.rwlock.Lock()
	defer db.rwlock.Unlock()

	db.opts = levigo.NewOptions()
	cache := levigo.NewLRUCache(3<<30)
	db.opts.SetCache(cache)
	db.opts.SetCreateIfMissing(true)
	filter := levigo.NewBloomFilter(10)
	db.opts.SetFilterPolicy(filter)

	db.ro = levigo.NewReadOptions()
	db.wo = levigo.NewWriteOptions()


	db.meta, db.recsep, err = openMeta(db.dbdir, db.ro, db.wo, db.opts)
	if err != nil {return err}

	db.elements, err = openElements(db.dbdir, db.opts)
	if err != nil {return err}

	db.edges, err = openEdges(db.dbdir, db.opts)
	if err != nil {return err}

	db.hexaindex, err = openHexaIndex(db.dbdir, db.opts)
	if err != nil {return err}

	db.props, err = openProps(db.dbdir, db.opts)
	if err != nil {return err}

	db.keepcount(VertexType, 0)
	db.keepcount(EdgeType, 0)
	db.IsOpen = true
	return nil
}


func OpenGraph(dbdir string) (*DBGraph, error ) {
	if dbdir == "" { return nil, NoDirectory }

	db := new(DBGraph)
	db.dbdir = dbdir
	db.IsOpen = false

	err := db.Open()
	if err != nil {return nil, err }
	return db, nil
}

func (db *DBGraph) Clear() (error) {
	dbdir := db.dbdir
	db.Close()
	os.RemoveAll(dbdir)
	return db.Open()
}


func (db *DBGraph) Close() (bool, error) {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	db.IsOpen = false
	db.meta.Close()
	db.elements.Close()
	db.edges.Close()
	db.hexaindex.Close()
	db.props.Close()
	db.opts.Close()
	db.ro.Close()
	db.wo.Close()
	db.meta = nil
	db.elements = nil
	db.edges = nil
	db.hexaindex = nil
	db.props = nil
	db.opts = nil
	db.ro = nil
	db.wo = nil
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
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
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
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	return db.delVertex(vertex)
}

func (db *DBGraph) delVertex(vertex *DBVertex) error {
	if vertex == nil {	return NilValue }
	if vertex.DBElement == nil { return NilValue }
	id := vertex.Id()
	if id == nil {	return NilValue }
	val,err := db.elements.Get(db.ro, id)
	if err != nil {return err}
	if val == nil {return KeyDoesNotExist}

	// delete all hexastore data
	for _, edge := range db.vertexEdges(0, vertex) {
		db.delEdge(edge)
	}

	for _, edge := range db.vertexEdges(1, vertex) {
		db.delEdge(edge)
	}

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

	// now delete the vertex
	err = db.elements.Delete(db.wo, id)
	if err != nil {return err}
	
	db.keepcount(VertexType, -1)

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

func (db *DBGraph) vertexEdges(outorin int, vertex *DBVertex) ([]*DBEdge) {
	edges := []*DBEdge{}
	if vertex == nil || vertex.id == nil { return edges }
	
	// outorin == 0 is out, 1 = in
	//prefix := 
	var prefix []byte
	if outorin == 0 {
		prefix = joinBytes(db.recsep, []byte(SPO), vertex.id)
	} else if outorin == 1 {
		prefix = joinBytes(db.recsep, []byte(OPS), vertex.id)
	} else {
		return edges
	}
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.hexaindex.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	it.Seek(prefix)


	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		//hxrec := it.Key()
		_, _, eid, _ := idsFromHexaKey(db.recsep, it.Key())
		//fmt.Printf("v=%v, eid=%v\n", string(vertex.Id()[:]),string(hxrec[:]))
		edges = append(edges, db.Edge(eid))
	}

	return edges

}

func (db *DBGraph) VertexOutEdges(vertex *DBVertex) ([]*DBEdge) {
	return db.vertexEdges(0, vertex)
}

func (db *DBGraph) VertexInEdges(vertex *DBVertex) ([]*DBEdge) {
	return db.vertexEdges(1, vertex)
}

func idsFromHexaKey(sep []byte, elements []byte) ([]byte, []byte, []byte, error) {
	if elements == nil { return nil, nil, nil, NilValue }
	elementarr := bytes.Split(elements, sep)
	if len(elementarr) != 4 {return nil,nil,nil, InvalidParameterValue }
	var subject, object, predicate []byte

	switch HexaIndexType(elementarr[0]) {
		case SPO:
			subject = elementarr[1]
			object = elementarr[3]
			predicate = elementarr[2]
		case SOP:
			subject = elementarr[1]
			object = elementarr[2]
			predicate = elementarr[3]
		case OPS:
			subject = elementarr[3]
			object = elementarr[1]
			predicate = elementarr[2]
		case OSP:
			subject = elementarr[2]
			object = elementarr[1]
			predicate = elementarr[3]
		case PSO:
			subject = elementarr[2]
			object = elementarr[3]
			predicate = elementarr[1]
		case POS:
			subject = elementarr[3]
			object = elementarr[2]
			predicate = elementarr[1]
		default:
			subject, object, predicate = nil, nil, nil
	}

	if subject == nil { return nil, nil, nil, InvalidParameterValue }
	return subject, object, predicate, nil

}

func newHexaIndexKey(sep []byte, subject []byte, object []byte, predicate []byte) (*HexaIndexKeys, error) {
	if sep == nil || subject == nil || object ==nil || predicate == nil { return nil, NilValue }
	hi := new(HexaIndexKeys)
	hi.spo = joinBytes(sep, []byte(SPO), subject, predicate, object)
	hi.sop = joinBytes(sep, []byte(SOP), subject, object, predicate)
	hi.ops = joinBytes(sep, []byte(OPS), object, predicate, subject)
	hi.osp = joinBytes(sep, []byte(OSP), object, subject, predicate)
	hi.pso = joinBytes(sep, []byte(PSO), predicate, subject, object)
	hi.pos = joinBytes(sep, []byte(POS), predicate, object, subject)
	return hi, nil
}


func (db *DBGraph) AddEdge(id []byte, outvertex *DBVertex, invertex *DBVertex, label string) (*DBEdge, error) {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
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

	labelbyte := []byte(label)
	//err = db.edges.Put(db.wo, id, db.toEdgeRecord(outvertex, invertex, label))
	err = db.edges.Put(db.wo, id, joinBytes(db.recsep, outvertex.id, invertex.id, labelbyte))
	if err != nil {return nil, err}

	// //todo - add hexascale index
	hi, _ := newHexaIndexKey(db.recsep, outvertex.id, invertex.id, id)

	wb := levigo.NewWriteBatch()
	defer wb.Close()
	wb.Put(hi.spo, labelbyte)
	wb.Put(hi.sop, labelbyte)
	wb.Put(hi.ops, labelbyte)
	wb.Put(hi.osp, labelbyte)
	wb.Put(hi.pso, labelbyte)
	wb.Put(hi.pos, labelbyte)

	err = db.hexaindex.Write(db.wo, wb)
	if err != nil {return nil, err}

	db.keepcount(EdgeType, 1)
	return edge, nil
}

func joinBytes(sep []byte, elements ...[]byte) ([]byte) {
	if len(elements) < 1 { return []byte{} } 
	return bytes.Join(elements, sep)
}

/*
func splitBytes(sep []byte, elements []byte) (int, [][]byte) {
	if elements == nil {return 0, nil}
	elementarr := bytes.Split(elements, sep)
	n := len(elementarr)
	return n, elementarr
}
*/

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
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	return db.delEdge(edge) 
}

func (db *DBGraph) delEdge(edge *DBEdge) error {
	if edge == nil {	return NilValue }
	if edge.DBElement == nil { return NilValue }
	id := edge.Id()
	if id == nil {	return NilValue }
	val,err := db.elements.Get(db.ro, id)
	if err != nil {return err}
	if val == nil {return KeyDoesNotExist}


	//  delete all hexastore data
	// //todo - add hexascale index
	//hexaindex key 
	hi, _ := newHexaIndexKey(db.recsep, edge.VertexOut().id, edge.VertexIn().id, id)
	wb2 := levigo.NewWriteBatch()
	defer wb2.Close()
	wb2.Delete(hi.spo)
	wb2.Delete(hi.sop)
	wb2.Delete(hi.ops)
	wb2.Delete(hi.osp)
	wb2.Delete(hi.pso)
	wb2.Delete(hi.pos)

	err = db.hexaindex.Write(db.wo, wb2)
	if err != nil {return err}

	// delete all properties data 
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.props.NewIterator(ro)
	defer ro.Close()
	defer it.Close()
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

	// now delete the edge
	err = db.elements.Delete(db.wo, id)
	if err != nil {return err}
	
	err = db.edges.Delete(db.wo, id)
	db.keepcount(EdgeType, -1)
	

	return nil
}

func (db *DBGraph) Edges() []*DBEdge {
	edges := []*DBEdge{}
	var edge *DBEdge
	ro := levigo.NewReadOptions()
	defer ro.Close()
	ro.SetFillCache(false)
	it := db.edges.NewIterator(ro)
	defer it.Close()
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

func (db *DBGraph) getPropKey(id []byte, prop string) ([]byte) {
	keyvalues := [][]byte{}
	keyvalues = append(keyvalues,id, []byte(prop))
	key := bytes.Join(keyvalues, db.recsep)
	return key
}

func (db *DBGraph) ElementProperty(element *DBElement, prop string) ([]byte) {
	if prop == "" || element == nil || element.id == nil { return nil }
	key := db.getPropKey(element.id, prop)
	val, err := db.props.Get(db.ro, key)
	if err != nil {return nil}
	return val
}

func (db *DBGraph) ElementSetProperty(element *DBElement, prop string, value []byte) (error){
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	if prop == "" || element == nil || element.id == nil { return nil }
	key := db.getPropKey(element.id, prop)
	err := db.props.Put(db.wo, key, value)
	return err
}

func (db *DBGraph) ElementDelProperty(element *DBElement, prop string) ([]byte) {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	if prop == "" || element == nil || element.id == nil { return nil }
	key := db.getPropKey(element.id, prop)
	val, err := db.props.Get(db.ro, key)
	if err != nil {return nil}
	err = db.props.Delete(db.wo, key)
	if err != nil {return nil}
	return val
}

func (db *DBGraph) ElementPropertyKeys(element *DBElement) ([]string) {
	propkeys := []string{}
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.props.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	prefix := append(element.id, db.recsep...)
	it.Seek(prefix)
	var prop []byte
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		prop = bytes.Split(it.Key(), db.recsep)[1]
		propkeys = append(propkeys, string(prop[:]))
	}
	return propkeys
}
