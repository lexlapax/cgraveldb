package levigo


import (
		"bytes"
		"fmt"
		"errors"
		"os"
		"sync"
		"encoding/binary"
		"strconv"
		"github.com/jmhodges/levigo"
		"path"
		"github.com/lexlapax/graveldb/core"
		"github.com/lexlapax/graveldb/util"
)

const (
	GraphImpl                = "levigo"
	)


// vars Errors in addition to IO errors.
var (
	NoDirectory     = errors.New("need to pass a valid path for a db directory")
	InvalidParameterValue = errors.New("the value of parameter passed is invalid")
	DBNotOpen = errors.New("db is not open")
	metadb = "meta.db"	
	atomdb = "atom.db"
	hexaindexdb = "hexaindex.db"
	propdb = "prop.db"
	recsep = "\x1e" //the ascii record separator dec val 30
	fieldsep = "\x1f" // the ascii unit separator dec val 31
	propVertiiCount = "vertiicount"
	propEdgeCount = "edgecount"
	propRecSep = "recsep"
	propNextId = "nextid"
	propUuid = "uuid"
	register sync.Once
)

func Register() {
	register.Do(func() {core.Register(GraphImpl, &GraphLevigo{})} )
}

func init() {
	Register()
}


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


type GraphLevigo struct {
	meta *levigo.DB
	atoms *levigo.DB
	hexaindex *levigo.DB
	props *levigo.DB
	dbdir string
	opts *levigo.Options
	ro *levigo.ReadOptions
	wo *levigo.WriteOptions
	keyindex *KeyIndex
	recsep []byte
	isopen bool
	rwlock *sync.RWMutex
	caps *graphCaps
	uuid string
}

func (db *GraphLevigo) nextId() string {
	lastidbyte, _ := db.meta.Get(db.ro, []byte(propNextId))
	if lastidbyte == nil { 
		lastidbyte = []byte(strconv.FormatUint(uint64(0), 16)) 
	}

	for {
		lastidstr := string(lastidbyte[:])
		intcounter, _ := strconv.ParseUint(lastidstr, 16, 64)
		intcounter++ 
		nextidstr := strconv.FormatUint(intcounter, 16)
		lastidbyte = []byte(nextidstr)
		val, _ := db.atoms.Get(db.ro, lastidbyte)
		if val == nil  {
			break
		}
	}
	db.meta.Put(db.wo, []byte(propNextId), lastidbyte)
	return string(lastidbyte[:])
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
	return levigo.Open(path.Join(dbdir, atomdb), opts)
}

func openHexaIndex(dbdir string, opts *levigo.Options) (*levigo.DB, error) {
	return levigo.Open(path.Join(dbdir, hexaindexdb), opts)
}

func openProps(dbdir string, opts *levigo.Options) (*levigo.DB, error) {
	return levigo.Open(path.Join(dbdir, propdb), opts)
}

func (db *GraphLevigo) Capabilities() core.GraphCaps {
	return db.caps
}


func (db *GraphLevigo) Open(args ...interface{}) error {
	if db.isopen == true {return nil }//errors.New("db already open") }
	if len(args) > 0 {
		if aString, found := args[0].(string); found {
			if aString == "" {
				return NoDirectory
			} else {
				db.dbdir = aString
			}
		} else {
			return InvalidParameterValue
		}
	}
	if db.dbdir == "" { return NoDirectory }

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

	db.caps = new(graphCaps)

	db.meta, db.recsep, err = openMeta(db.dbdir, db.ro, db.wo, db.opts)
	if err != nil {return err}

	db.atoms, err = openElements(db.dbdir, db.opts)
	if err != nil {return err}


	db.hexaindex, err = openHexaIndex(db.dbdir, db.opts)
	if err != nil {return err}

	db.props, err = openProps(db.dbdir, db.opts)
	if err != nil {return err}

	db.keepcount(core.VertexType, 0)
	db.keepcount(core.EdgeType, 0)
	if db.caps.KeyIndex() {
		db.keyindex = NewKeyIndex(db.dbdir)
	}
	uuidbytes, err := db.getDbProperty(propUuid)
	if uuidbytes == nil {
		uuid, err := util.UUID()
		if err != nil { return err }
		uuidbytes = []byte(uuid)
		db.putDbProperty(propUuid, uuidbytes)
		db.uuid = uuid
	} else {
		db.uuid = string(uuidbytes[:])
	}
	db.isopen = true
	return nil
}

func OpenGraph(dbdir string) (*GraphLevigo, error ) {
	if dbdir == "" { return nil, NoDirectory }

	db := new(GraphLevigo)
	db.dbdir = dbdir
	db.isopen = false

	err := db.Open()
	if err != nil {return nil, err }
	return db, nil
}

func (db *GraphLevigo) Clear() error {
	dbdir := db.dbdir
	db.Close()
	os.RemoveAll(dbdir)
	return db.Open()
}

func (db *GraphLevigo) Close() error {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	db.isopen = false
	db.meta.Close()
	db.atoms.Close()
	db.hexaindex.Close()
	db.props.Close()
	db.opts.Close()
	db.ro.Close()
	db.wo.Close()
	if db.caps.KeyIndex() {
		db.keyindex.Close()
	}
	db.meta = nil
	db.atoms = nil
	db.hexaindex = nil
	db.props = nil
	db.opts = nil
	db.ro = nil
	db.wo = nil
	db.keyindex = nil
	return nil
}


func (db *GraphLevigo) IsOpen() bool {
	return db.isopen
}

func (db *GraphLevigo) String() (string) {
	str := fmt.Sprintf("#GraphLevigo:dbdir=%v#",db.dbdir)
	return str
}

func(db *GraphLevigo) getDbProperty(prop string) ([]byte, error){
	if prop == "" {return nil, core.ErrNilValue}
	val, err := db.meta.Get(db.ro, []byte(prop))
	if err != nil {return nil, err}
	return val, nil
}

func(db *GraphLevigo) putDbProperty(prop string, val []byte) ([]byte, error){
	if prop == "" {return nil, core.ErrNilValue}
	key := []byte(prop)
	oldval, err := db.meta.Get(db.ro, key)
	if err != nil {return nil, err}
	err2 := db.meta.Put(db.wo, key, val)
	if err2 != nil {return nil, err}
	return oldval, nil
}

func (db *GraphLevigo) keepcount(etype core.AtomType, upordown int) (uint) {
	var storedcount, returncount uint
	var key string
	if etype == core.VertexType {
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
		tempint,_ := binary.Uvarint(val)
		storedcount = uint(tempint)
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
		bufsize := binary.Size(uint64(returncount))
		buf := make([]byte, bufsize)
		binary.PutUvarint(buf, uint64(returncount))
		db.putDbProperty(key, buf)
	}
	return returncount
}

func (graph *GraphLevigo) Guid() string {
	return graph.uuid	
}

func (db *GraphLevigo) AddVertex(id string) (core.Vertex, error) {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	if id == "" {
		id = db.nextId()
	} else {
		val, _ := db.atoms.Get(db.ro, []byte(id))
		if val != nil {
			return nil, core.ErrAlreadyExists
		}
	}
	err := db.atoms.Put(db.wo, []byte(id), []byte(core.VertexType))
	if err != nil {return nil, err}
	vertex := &VertexLevigo{&AtomLevigo{db,id,core.VertexType}}
	db.keepcount(core.VertexType, 1)
	return vertex, nil
}

func (db *GraphLevigo) Vertex(id string) (core.Vertex, error) {
	if id == "" {return nil,core.ErrNilValue}
	val,err := db.atoms.Get(db.ro, []byte(id))
	if err != nil {return nil, err}
	if val == nil {return nil, nil}
	if core.AtomType(val) != core.VertexType {return nil, nil}
	vertex := &VertexLevigo{&AtomLevigo{db, id, core.VertexType}}
	return vertex, nil
}

func (db *GraphLevigo) DelVertex(vertex core.Vertex) error {
	if vertex == nil {	return core.ErrNilValue }
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	vertexlevigo := vertex.(*VertexLevigo)
	return db.delVertex(vertexlevigo)
}

func (db *GraphLevigo) delVertex(vertex *VertexLevigo) error {
	if vertex == nil {	return core.ErrNilValue }
	if vertex.AtomLevigo == nil { return core.ErrNilValue }
	id := vertex.Id()
	if id == "" {	return core.ErrNilValue }
	val,err := db.atoms.Get(db.ro, []byte(id))
	if err != nil {return err}
	if val == nil {return nil}

	// delete all hexastore data
	edges, _ :=  db.vertexEdges(core.DirOut, vertex) 
	for _, edge := range edges {
		db.delEdge(edge.(*EdgeLevigo))
	}

	edges, _ =  db.vertexEdges(core.DirIn, vertex) 
	for _, edge := range edges {
		db.delEdge(edge.(*EdgeLevigo))
	}

	// delete all properties data 
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.props.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	prefix := append([]byte(id), db.recsep...)
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
	err = db.atoms.Delete(db.wo, []byte(id))
	if err != nil {return err}
	
	db.keepcount(core.VertexType, -1)

	return nil
}

func (db *GraphLevigo) Vertices() ([]core.Vertex, error) {
	vertii := []core.Vertex{}
	var vertex *VertexLevigo
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.atoms.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	it.SeekToFirst()
	for it = it; it.Valid(); it.Next() {
		if core.AtomType(it.Value()) == core.VertexType {
			vertex = &VertexLevigo{&AtomLevigo{db, string(it.Key()[:]), core.VertexType}}
			vertii = append(vertii, vertex)
		}
	}
	return vertii, nil
}


func (db *GraphLevigo) vertexVertices(outorin core.Direction, vertex *VertexLevigo, labels ...string) ([]core.Vertex, error) {
	vertices := []core.Vertex{}
	if vertex == nil || vertex.id == "" { return vertices, core.ErrNilValue }
	
	// outorin == 0 is out, 1 = in
	//prefix := 
	var prefix []byte
	if outorin == core.DirOut {
		prefix = joinBytes(db.recsep, []byte(SPO), []byte(vertex.Id()))
	} else if outorin == core.DirIn {
		prefix = joinBytes(db.recsep, []byte(OPS), []byte(vertex.Id()))
	} else {
		return vertices, core.ErrDirAnyUnsupported
	}
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.hexaindex.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	it.Seek(prefix)
	labelset := util.NewStringSet()
	//fmt.Printf("labels = %v\n", labels)
	if len(labels) > 0 {
		for _, label := range labels {
			labelset.Add(label)
		}
	}

	//fmt.Printf("labelset=%v\n", labelset)

	addvertex := false
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		//hxrec := it.Key()
		outvid, invid, _, _ := idsFromHexaKey(db.recsep, it.Key())
		if labelset.Count() > 0 {
			label := string(it.Value()[:])
			if labelset.Contains(label) { addvertex = true }
		} else {
			addvertex = true
		}
		if addvertex == true {
			gotvertexid := []byte{}
			//which vertex
			if bytes.Compare([]byte(vertex.Id()), outvid) == 0 {
				gotvertexid = invid
			} else {
				gotvertexid = outvid
			}
			//fmt.Printf("v=%v, eid=%v\n", string(vertex.Id()[:]),string(hxrec[:]))
			gotvertex, _ := db.Vertex(string(gotvertexid[:]))
			vertices = append(vertices, gotvertex)
		}
		addvertex = false
	}

	return vertices, nil

}

func (db *GraphLevigo) vertexEdges(outorin core.Direction, vertex *VertexLevigo, labels ...string) ([]core.Edge, error) {
	edges := []core.Edge{}
	if vertex == nil || vertex.id == "" { return edges, core.ErrNilValue }
	
	// outorin == 0 is out, 1 = in
	//prefix := 
	var prefix []byte
	if outorin == core.DirOut {
		prefix = joinBytes(db.recsep, []byte(SPO), []byte(vertex.Id()))
	} else if outorin == core.DirIn {
		prefix = joinBytes(db.recsep, []byte(OPS), []byte(vertex.Id()))
	} else {
		return edges, core.ErrDirAnyUnsupported
	}
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.hexaindex.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	it.Seek(prefix)
	labelset := util.NewStringSet()
	//fmt.Printf("labels = %v\n", labels)
	if len(labels) > 0 {
		for _, label := range labels {
			labelset.Add(label)
		}
	}

	//fmt.Printf("labelset=%v\n", labelset)

	addedge := false
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		//hxrec := it.Key()
		_, _, eid, _ := idsFromHexaKey(db.recsep, it.Key())
		if labelset.Count() > 0 {
			label := string(it.Value()[:])
			if labelset.Contains(label) { addedge = true }
		} else {
			addedge = true
		}
		if addedge == true {
			//fmt.Printf("v=%v, eid=%v\n", string(vertex.Id()[:]),string(hxrec[:]))
			edge, _ := db.Edge(string(eid[:]))
			edges = append(edges, edge)
		}
		addedge = false
	}

	return edges, nil

}

func (db *GraphLevigo) VertexOutEdges(vertex *VertexLevigo, labels ...string) ([]core.Edge, error) {
	return db.vertexEdges(core.DirOut, vertex, labels...)
}

func (db *GraphLevigo) VertexInEdges(vertex *VertexLevigo, labels ...string) ([]core.Edge, error) {
	return db.vertexEdges(core.DirIn, vertex, labels...)
}

// returns outvertex, invertex, edge,  error
func idsFromHexaKey(sep []byte, atoms []byte) ([]byte, []byte, []byte, error) {
	if atoms == nil { return nil, nil, nil, core.ErrNilValue }
	atomarr := bytes.Split(atoms, sep)
	if len(atomarr) != 4 {return nil,nil,nil, InvalidParameterValue }
	var subject, object, predicate []byte

	switch HexaIndexType(atomarr[0]) {
		case SPO:
			subject = atomarr[1]
			object = atomarr[3]
			predicate = atomarr[2]
		case SOP:
			subject = atomarr[1]
			object = atomarr[2]
			predicate = atomarr[3]
		case OPS:
			subject = atomarr[3]
			object = atomarr[1]
			predicate = atomarr[2]
		case OSP:
			subject = atomarr[2]
			object = atomarr[1]
			predicate = atomarr[3]
		case PSO:
			subject = atomarr[2]
			object = atomarr[3]
			predicate = atomarr[1]
		case POS:
			subject = atomarr[3]
			object = atomarr[2]
			predicate = atomarr[1]
		default:
			subject, object, predicate = nil, nil, nil
	}

	if subject == nil { return nil, nil, nil, InvalidParameterValue }
	return subject, object, predicate, nil

}

func newHexaIndexKey(sep []byte, subject string, object string, predicate string) (*HexaIndexKeys, error) {
	if sep == nil || subject == "" || object == "" || predicate == "" { return nil, core.ErrNilValue }
	hi := new(HexaIndexKeys)
	hi.spo = joinBytes(sep, []byte(SPO), []byte(subject), []byte(predicate), []byte(object))
	hi.sop = joinBytes(sep, []byte(SOP), []byte(subject), []byte(object), []byte(predicate))
	hi.ops = joinBytes(sep, []byte(OPS), []byte(object), []byte(predicate), []byte(subject))
	hi.osp = joinBytes(sep, []byte(OSP), []byte(object), []byte(subject), []byte(predicate))
	hi.pso = joinBytes(sep, []byte(PSO), []byte(predicate), []byte(subject), []byte(object))
	hi.pos = joinBytes(sep, []byte(POS), []byte(predicate), []byte(object), []byte(subject))
	return hi, nil
}


func (db *GraphLevigo) AddEdge(id string, outvertex core.Vertex, invertex core.Vertex, label string) (core.Edge, error) {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	if (outvertex == nil) {return nil, core.ErrNilValue}
	if (invertex == nil) {return nil, core.ErrNilValue}
	if id == "" {
		id = db.nextId()
	} else {
		val, _ := db.atoms.Get(db.ro, []byte(id))
		if val != nil {
			return nil, core.ErrAlreadyExists
		}
	}
	idbyte := []byte(id)
	err := db.atoms.Put(db.wo, idbyte, []byte(core.EdgeType))
	if err != nil {return nil, err}
	edge := &EdgeLevigo{&AtomLevigo{db, id, core.EdgeType}, outvertex.(*VertexLevigo), invertex.(*VertexLevigo), label}

	labelbyte := []byte(label)

	// - add hexascale index
	hi, _ := newHexaIndexKey(db.recsep, outvertex.Id(), invertex.Id(), id)

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

	db.keepcount(core.EdgeType, 1)
	return edge, nil
}

func joinBytes(sep []byte, atoms ...[]byte) ([]byte) {
	if len(atoms) < 1 { return []byte{} } 
	return bytes.Join(atoms, sep)
}

func (db *GraphLevigo) fromEdgeRecord(record []byte) (core.Vertex, core.Vertex, string) {
	if record == nil { return nil, nil, ""}
	edgevalues := bytes.Split(record, db.recsep)

	outvertex,_ := db.Vertex(string(edgevalues[0]))
	invertex,_ := db.Vertex(string(edgevalues[1]))
	label := string(edgevalues[2][:])
	return outvertex, invertex, label
}

func (db *GraphLevigo) Edge(id string) (core.Edge, error) {
	if id == "" {return nil, core.ErrNilValue}
	val,err := db.atoms.Get(db.ro, []byte(id))
	if err != nil {return nil, err}
	if val == nil {return nil, core.ErrNilValue}
	if core.AtomType(val) != core.EdgeType {return nil, nil}

	prefix := joinBytes(db.recsep, []byte(PSO), []byte(id))
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.hexaindex.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	it.Seek(prefix)
	if it.Valid() && bytes.HasPrefix(it.Key(), prefix) {
		outvertexid, invertexid, eid, _  := idsFromHexaKey(db.recsep, it.Key())
		outvertex := &VertexLevigo{&AtomLevigo{db, string(outvertexid[:]), core.VertexType}}
		invertex := &VertexLevigo{&AtomLevigo{db, string(invertexid[:]), core.VertexType}}
		edge := &EdgeLevigo{&AtomLevigo{db, string(eid[:]), core.EdgeType}, outvertex, invertex, string(it.Value()[:])}
		return edge, nil
	} else {
		return nil, nil
	}
}


func (db *GraphLevigo) DelEdge(edge core.Edge) error {
	if edge == nil {	return core.ErrNilValue }
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	edgelevigo := edge.(*EdgeLevigo)
	return db.delEdge(edgelevigo)
}

func (db *GraphLevigo) delEdge(edge *EdgeLevigo) error {
	if edge == nil {	return core.ErrNilValue }
	if edge.AtomLevigo == nil { return core.ErrNilValue }
	id := edge.Id()
	if id == "" {	return core.ErrNilValue }
	val,err := db.atoms.Get(db.ro, []byte(id))
	if err != nil {return err}
	if val == nil {return core.ErrDoesntExist}


	//  delete all hexastore data
	// //todo - add hexascale index
	//hexaindex key 
	vout, _ := edge.VertexOut()
	vin, _ := edge.VertexIn()
	hi, _ := newHexaIndexKey(db.recsep, vout.Id(), vin.Id(), id)
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
	prefix := append([]byte(id), db.recsep...)
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
	err = db.atoms.Delete(db.wo, []byte(id))
	if err != nil {return err}
	
	db.keepcount(core.EdgeType, -1)
	

	return nil
}

func (db *GraphLevigo) Edges() ([]core.Edge, error) {
	edges := []core.Edge{}

	prefix := []byte(PSO)
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.hexaindex.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	it.Seek(prefix)

	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		outvertexid, invertexid, eid, _  := idsFromHexaKey(db.recsep, it.Key())
		outvertex := &VertexLevigo{&AtomLevigo{db, string(outvertexid[:]), core.VertexType}}
		invertex := &VertexLevigo{&AtomLevigo{db, string(invertexid[:]), core.VertexType}}
		edge := &EdgeLevigo{&AtomLevigo{db, string(eid[:]), core.EdgeType}, outvertex, invertex, string(it.Value()[:])}
		edges = append(edges, edge)
	}

	return edges, nil

}


func (db *GraphLevigo) EdgeCount() uint {
	return db.keepcount(core.EdgeType, 0)
}

func (db *GraphLevigo) VertexCount() uint {
	return db.keepcount(core.VertexType, 0)
}

func (db *GraphLevigo) getPropKey(id []byte, prop string) ([]byte) {
	keyvalues := [][]byte{}
	keyvalues = append(keyvalues,id, []byte(prop))
	key := bytes.Join(keyvalues, db.recsep)
	return key
}

func (db *GraphLevigo) AtomProperty(atom *AtomLevigo, prop string) ([]byte, error) {
	if prop == "" || atom == nil || atom.id == "" { return nil, nil }
	key := db.getPropKey([]byte(atom.id), prop)
	val, err := db.props.Get(db.ro, key)
	if err != nil {return nil, err}
	return val, nil
}

func (db *GraphLevigo) AtomSetProperty(atom *AtomLevigo, prop string, value []byte) error {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	if prop == "" || atom == nil || atom.id == "" { return nil }
	key := db.getPropKey([]byte(atom.id), prop)
	oldvalue, _ := db.props.Get(db.ro, key)
	err := db.props.Put(db.wo, key, value)

	if db.caps.KeyIndex() {
		db.keyindex.update(prop, string(value[:]), string(oldvalue[:]), atom.id, atom.atomType)
	}

	return err
}

func (db *GraphLevigo) AtomDelProperty(atom *AtomLevigo, prop string) error {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	if prop == "" || atom == nil || atom.id == "" { return nil }
	key := db.getPropKey([]byte(atom.id), prop)
	//val, err := db.props.Get(db.ro, key)
	//if err != nil {return err}
	oldvalue, _ := db.props.Get(db.ro, key)

	err := db.props.Delete(db.wo, key)

	if db.caps.KeyIndex() {
		db.keyindex.remove(prop, string(oldvalue[:]), atom.id, atom.atomType)
	}

	return err
}

func (db *GraphLevigo) AtomPropertyKeys(atom *AtomLevigo) ([]string, error) {
	propkeys := []string{}
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.props.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	prefix := append([]byte(atom.id), db.recsep...)
	it.Seek(prefix)
	var prop []byte
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		prop = bytes.Split(it.Key(), db.recsep)[1]
		propkeys = append(propkeys, string(prop[:]))
	}
	return propkeys, nil
}


func (db *GraphLevigo) CreateKeyIndex(key string, atomType core.AtomType) error {
	return db.keyindex.createKeyIndex(key, atomType)
}

func (db *GraphLevigo) DropKeyIndex(key string, atomType core.AtomType) error {
	return db.keyindex.dropKeyIndex(key, atomType)
}

func (db *GraphLevigo) IndexedKeys(atomType core.AtomType) []string {
		return db.keyindex.indexedKeys(atomType)
}

func (db *GraphLevigo) VerticesWithProp(key string, value string) []core.Vertex {
	ids := db.keyindex.searchIds(key, value, core.VertexType)
	vertices := []core.Vertex{}
	var vertex core.Vertex
	for _, idstring := range ids {
		vertex, _ = db.Vertex(idstring)
		if vertex != nil {
			vertices = append(vertices, vertex)
		}
	}
	return vertices
}

func (db *GraphLevigo) EdgesWithProp(key string, value string) []core.Edge {
	ids := db.keyindex.searchIds(key, value, core.EdgeType)
	edges := []core.Edge{}
	var edge core.Edge
	for _, idstring := range ids {
		edge, _ = db.Edge(idstring)
		if edge != nil {
			edges = append(edges, edge)
		}
	}
	return edges
}
