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

func TestOpenGraph(t *testing.T){
	//t.Skip()
	gdb, dberr := OpenGraph(dbdir)
	if dberr != nil { t.Fatal(dberr) }

	if gdb == nil {
		t.Error("graphdb should not be nil")
	} else {
		gdb.Clear() 
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
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

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
}


func TestCloseAndOpen(t *testing.T) {
	//t.Skip()
	gdb, dberr := OpenGraph(dbdir)

	if dberr != nil { t.Fatal(dberr) }
	gdb.Clear()
	if gdb.EdgeCount() != 0 { t.Error("should have 0 edges")}
	if gdb.VertexCount() != 0 { t.Error("should have 0 vertices")}
	gdb.AddVertex([]byte("somerandomstringid"))
	if gdb.VertexCount() != 1 { t.Error("should have 1 vertex")}
	gdb.Close()
	gdb, dberr = OpenGraph(dbdir)
	if gdb.VertexCount() != 1 { t.Error("should have 1 vertex")}
	gdb.Close()
}


func TestGetVertex(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

	ida := []byte("somerandomstringid")
	idb := []byte("idb")
	vertexa, _ := gdb.AddVertex(ida)
	vertexb  := gdb.Vertex(ida)
	assert.Equal(t, vertexa, vertexb)
	assert.True(t, gdb.Vertex(idb) == nil)
	vertexc, _ := gdb.AddVertex(idb)
	vertexd := gdb.Vertex(idb)
	assert.Equal(t, vertexc, vertexd)
}

func TestDelVertex(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

	ida := []byte("somerandomstringid")

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
}

func TestVertexCount(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

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
}


func TestVertices(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

	ida := []byte("somerandomstringid")
	testvertii := []*DBVertex{}
	var vertex *DBVertex

	assert.True(t, len(gdb.Vertices()) == 0)

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
}

func TestAddEdge(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

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
		assert.Equal(t, vertex1, edge1.VertexOut())
		assert.Equal(t, vertex2, edge1.VertexIn())
		assert.Equal(t, "edgeforward", edge1.Label())
	}

	edge2, errb := gdb.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.True(t, edge2 == nil)
	assert.Equal(t, KeyExists, errb )
}

func TestGetEdge(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")
	assert.True(t, gdb.Edge(eid1) == nil)

	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)
	edge1, _ := gdb.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	edge1a := gdb.Edge(eid1)
	assert.Equal(t, edge1, edge1a)

	//allow duplicates 
	eid2 := []byte("thisisedge2")
	edge2, _ := gdb.AddEdge(eid2, vertex1, vertex2, "edgeforward")
	assert.Equal(t, eid2, edge2.Id())
	assert.Equal(t, vertex1, edge1.VertexOut())
	assert.Equal(t, vertex2, edge1.VertexIn())
	assert.Equal(t, "edgeforward", edge1.Label())
}

func TestDelEdge(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")

	err := gdb.DelEdge(nil)
	assert.Equal(t, NilValue, err)
	err = gdb.DelEdge(new(DBEdge))
	assert.Equal(t, NilValue, err)


	edgenull := &DBEdge{new(DBElement),nil,nil,""}
	edgenull.db = gdb 
	edgenull.id = eid1 
	edgenull.Elementtype = EdgeType
	err = gdb.DelEdge(edgenull)
	assert.Equal(t, KeyDoesNotExist, err)

	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)
	edge1, _ := gdb.AddEdge(eid1, vertex1, vertex2, "edgeforward")

	err = gdb.DelEdge(edge1)
	assert.True(t, err == nil)
	assert.True(t, gdb.Edge(eid1) == nil)
	err = gdb.DelEdge(edge1)
	assert.Equal(t, KeyDoesNotExist, err)
}


func TestEdges(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")

	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)

	assert.True(t, len(gdb.Edges()) == 0)
	testedges := []*DBEdge{}
	var edge *DBEdge
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			edge, _= gdb.AddEdge([]byte(a + "-" + n), vertex1, vertex2, "somedge")
			testedges = append(testedges, edge)
		}
	}
	assert.Equal(t, testedges, gdb.Edges())
	lastedge, _ := gdb.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.NotEqual(t, testedges, gdb.Edges())
	//keys are lexicaly ordered.. lastedge should be the last in the list
	testedges = gdb.Edges()
	assert.Equal(t, lastedge, testedges[len(testedges) - 1])
}


func TestEdgeCount(t *testing.T) {
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

	assert.Equal(t, uint64(0), gdb.EdgeCount())

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")

	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)

	edge1,_ := gdb.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.Equal(t, uint64(1), gdb.EdgeCount())

	testedges := []*DBEdge{}
	var edge *DBEdge
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			edge, _= gdb.AddEdge([]byte(a + "-" + n), vertex1, vertex2, "somedge")
			testedges = append(testedges, edge)
		}
	}
	numv := len(testedges)
	assert.Equal(t, uint64(numv + 1), gdb.EdgeCount())
	gdb.DelEdge(edge1)
	assert.Equal(t, uint64(numv), gdb.EdgeCount())
	for i :=0; i < numv; i++ {
		gdb.DelEdge(testedges[i])
		assert.Equal(t, uint64(numv - (i + 1)), gdb.EdgeCount() )
	}
	assert.Equal(t, uint64(0), gdb.EdgeCount())
}

