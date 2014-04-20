package levelgraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
	"reflect"
	"bytes"
	//"fmt"
	//"github.com/jmhodges/levigo"
	//"github.com/lexlapax/graveldb/core"	
)


var (
	dbdir = "./testing.db"
)


func cleanup(dbdir string) {
	os.RemoveAll(dbdir)
}

func TestOpenGraph(t *testing.T){
	cleanup(dbdir)
	defer cleanup(dbdir)

	gdb, dberr := opengraph(dbdir)
	if dberr != nil { t.Fatal(dberr) }

	if gdb == nil {
		t.Error("graphdb should not be nil")
	} else { 
		if dbdir != gdb.dbdir { t.Error("dbdir not equal")}
		if reflect.TypeOf(gdb.meta).String() != "*levigo.DB" { t.Error("gdb not valid type")}
		if gdb.meta == nil { t.Error("meta is nil") }
		if gdb.elements == nil { t.Error("elements is nil") }
		if gdb.edges == nil { t.Error("edges is nil") }
		if gdb.hexaindex == nil { t.Error("hexaindex is nil") }
		if gdb.props == nil { t.Error("props is nil") }
		if bytes.Compare(gdb.recsep, []byte("\x1f")) != 0 { t.Error("recsep does not match") }
		if gdb.EdgeCount() != 0 { t.Error("should have 0 edges")}
		if gdb.VertexCount() != 0 { t.Errorf("should have 0 vertices has %", gdb.VertexCount())}

		fi, _ := os.Lstat(dbdir)
		if !fi.IsDir() { t.Error("dbdir should be a directory") }
		if fi.Name() != "testing.db" { t.Error("dbdir name should match") }
		if gdb.String() != "#DBGraph:dbdir=./testing.db#" { t.Error("String method does not match")}
		gdb.Close()
	}
}

func TestAddVertex(t *testing.T) {
	//t.Skip()

	cleanup(dbdir)
	defer cleanup(dbdir)
	gdb,_ := opengraph(dbdir)

	id := []byte("somerandomstringid")
	vertex, err := gdb.AddVertex(nil)
	assert.True(t, vertex == nil)
	assert.Equal(t, NilValue, err)

	vertex, err = gdb.AddVertex(id)
	if assert.True(t, vertex != nil) {
		assert.Equal(t, id, vertex.Id())
		assert.Equal(t, VertexType, vertex.Elementtype)
		assert.Equal(t, nil, err)
		assert.Equal(t, "<DBVertex:somerandomstringid@#DBGraph:dbdir=./testing.db#>", vertex.String())
	}
	vertexb, errb := gdb.AddVertex(id) 
	assert.True(t, vertexb == nil)
	assert.Equal(t, KeyExists, errb )

	gdb.Close()
}


func TestCloseAndOpen(t *testing.T) {
	//t.Skip()

	cleanup(dbdir)
	defer cleanup(dbdir)

	gdb, dberr := opengraph(dbdir)
	if dberr != nil { t.Fatal(dberr) }
	if gdb.EdgeCount() != 0 { t.Error("should have 0 edges")}
	if gdb.VertexCount() != 0 { t.Error("should have 0 vertices")}
	gdb.AddVertex([]byte("somerandomstringid"))
	if gdb.VertexCount() != 1 { t.Error("should have 1 vertex")}
	gdb.Close()
	gdb, dberr = opengraph(dbdir)
	if gdb.VertexCount() != 1 { t.Error("should have 1 vertex")}
	gdb.Close()
}


func TestGetVertex(t *testing.T) {
	//t.Skip()

	cleanup(dbdir)
	defer cleanup(dbdir)
	gdb,_ := opengraph(dbdir)

	ida := []byte("somerandomstringid")
	idb := []byte("idb")
	vertexa, _ := gdb.AddVertex(ida)
	vertexb  := gdb.Vertex(ida)
	assert.Equal(t, vertexa, vertexb)
	assert.True(t, gdb.Vertex(idb) == nil)
	vertexc, _ := gdb.AddVertex(idb)
	vertexd := gdb.Vertex(idb)
	assert.Equal(t, vertexc, vertexd)
	gdb.Close()
}

func TestDelVertex(t *testing.T) {
	//t.Skip()

	cleanup(dbdir)
	defer cleanup(dbdir)

	ida := []byte("somerandomstringid")
	gdb,_ := opengraph(dbdir)

	err := gdb.DelVertex(nil)
	assert.Equal(t, NilValue, err)
	err = gdb.DelVertex(new(DBVertex))
	assert.Equal(t, NilValue, err)

	vertexnull := &DBVertex{new(DBElement)}
	vertexnull.db = gdb 
	vertexnull.id = ida 
	vertexnull.Elementtype = VertexType
	err = gdb.DelVertex(vertexnull)
	assert.Equal(t, KeyDoesNotExist, err)
	vertexa, _ := gdb.AddVertex(ida)
	assert.Equal(t, vertexa, gdb.Vertex(ida))
	err = gdb.DelVertex(vertexa)
	assert.True(t, err == nil)
	assert.True(t, gdb.Vertex(ida) == nil)
	err = gdb.DelVertex(vertexa)
	assert.Equal(t, KeyDoesNotExist, err)
	gdb.Close()
}

func TestVertexCount(t *testing.T) {
	//t.Skip()

	cleanup(dbdir)
	defer cleanup(dbdir)

	gdb,_ := opengraph(dbdir)

	assert.Equal(t, uint64(0), gdb.VertexCount())

	ida := []byte("somerandomstringid")
	vertexa,_ := gdb.AddVertex(ida)
	assert.Equal(t, uint64(1), gdb.VertexCount())

	testvertii := []*DBVertex{}
	var vertex *DBVertex
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _ = gdb.AddVertex([]byte(a + "-" + n))
			testvertii = append(testvertii, vertex)
		}
	}
	numv := len(testvertii)
	assert.Equal(t, uint64(numv + 1), gdb.VertexCount())
	gdb.DelVertex(vertexa)
	assert.Equal(t, uint64(numv), gdb.VertexCount())
	for i :=0; i < numv; i++ {
		gdb.DelVertex(testvertii[i])
		assert.Equal(t, uint64(numv - (i + 1)), gdb.VertexCount() )
	}
	assert.Equal(t, uint64(0), gdb.VertexCount())
	gdb.Close()
}


func TestVertices(t *testing.T) {
	//t.Skip()

	cleanup(dbdir)
	defer cleanup(dbdir)

	ida := []byte("somerandomstringid")
	gdb,_ := opengraph(dbdir)
	testvertii := []*DBVertex{}
	var vertex *DBVertex
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _= gdb.AddVertex([]byte(a + "-" + n))
			testvertii = append(testvertii, vertex)
		}
	}
	assert.Equal(t, testvertii, gdb.Vertices())
	lastvertex, _ := gdb.AddVertex(ida)
	assert.NotEqual(t, testvertii, gdb.Vertices())
	//keys are lexicaly ordered.. lastvertex should be the last in the list
	testvertii = gdb.Vertices()
	assert.Equal(t, lastvertex, testvertii[len(testvertii) - 1])
	gdb.Close()
}

func TestAddEdge(t *testing.T) {
	//t.Skip()

	cleanup(dbdir)
	defer cleanup(dbdir)

	gdb,_ := opengraph(dbdir)
	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")
	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)
	//fmt.Printf("v=%v\n",vertex2)

	edge1, err := gdb.AddEdge(nil,vertex1,vertex2,"edgeforward")
	assert.True(t, edge1 == nil)
	assert.Equal(t, NilValue, err)

	edge1, err = gdb.AddEdge(eid1,nil,vertex2,"edgeforward")
	assert.True(t, edge1 == nil)
	assert.Equal(t, NilValue, err)

	edge1, err = gdb.AddEdge(eid1,vertex1,nil,"edgeforward")
	assert.True(t, edge1 == nil)
	assert.Equal(t, NilValue, err)

	edge1, err = gdb.AddEdge(eid1, vertex1, vertex2, "edgeforward")

	if assert.True(t, edge1 != nil) {
		assert.Equal(t, eid1, edge1.Id())
		assert.Equal(t, EdgeType, edge1.Elementtype)
		assert.Equal(t, nil, err)
		assert.Equal(t, "<DBEdge:thisisedge1,s=thisisvertex1,o=thisisvertex2,l=edgeforward@#DBGraph:dbdir=./testing.db#>", edge1.String())
	}

	edge2, errb := gdb.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.True(t, edge2 == nil)
	assert.Equal(t, KeyExists, errb )

	gdb.Close()
}

func TestGetEdge(t *testing.T) {
	//t.Skip()

	cleanup(dbdir)
	defer cleanup(dbdir)

	gdb,_ := opengraph(dbdir)
	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")
	assert.True(t, gdb.Edge(eid1) == nil)

	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)
	edge1, _ := gdb.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	edge2 := gdb.Edge(eid1)
	assert.Equal(t, edge1, edge2) 
	/*
	ida := []byte("somerandomstringid")
	idb := []byte("idb")
	vertexa, _ := gdb.AddVertex(ida)
	vertexb  := gdb.Vertex(ida)
	assert.Equal(t, vertexa, vertexb)
	assert.True(t, gdb.Vertex(idb) == nil)
	vertexc, _ := gdb.AddVertex(idb)
	vertexd := gdb.Vertex(idb)
	assert.Equal(t, vertexc, vertexd)
	*/
	gdb.Close()
}

