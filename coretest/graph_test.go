package coretest

import (
	"github.com/lexlapax/graveldb/core"
	"github.com/stretchr/testify/assert"
	"testing"
	//"os"
	//"reflect"
	//"github.com/jmhodges/levigo"
	//"fmt"
	"github.com/lexlapax/graveldb/mem"	
)

var graphimp = "mem"
func init() {
	mem.Register()
}

func TestGraphEmpty(t *testing.T){
	graph := core.GetGraph(graphimp)
	graph.Open()
	defer graph.Close()

	//todo - graph capabilites based 
	vertices, _ := graph.Vertices()
	assert.Equal(t, 0, len(vertices))
	edges, _ := graph.Edges()
	assert.Equal(t, 0, len(edges))
}


func TestGraphAdd(t *testing.T) {
	graph := core.GetGraph(graphimp)
	graph.Open()
	defer graph.Close()

	a,_ := graph.AddVertex(nil)
	b,_ := graph.AddVertex(nil)
    edge, _ := graph.AddEdge(nil, a, b, "knows")

	//todo - graph capabilites based 
	vertices, _ := graph.Vertices()
	assert.Equal(t, 2, len(vertices))
	edges, _ := graph.Edges()
	assert.Equal(t, 1, len(edges))

	graph.DelVertex(a)

	vertices, _ = graph.Vertices()
	assert.Equal(t, 1, len(vertices))
	edges, _ = graph.Edges()
	assert.Equal(t, 0, len(edges))
	err := graph.DelEdge(edge)
	assert.True(t, err != nil)
}

func TestGraphSetProperties(t *testing.T) {
	graph := core.GetGraph(graphimp)
	graph.Open()
	defer graph.Close()
	//todo graph capabilities based
	a,_ := graph.AddVertex(nil)
	b,_ := graph.AddVertex(nil)
	graph.AddEdge(nil, a, b, "knows")
	graph.AddEdge(nil, a, b, "knows")
	edges, _ := b.Edges(core.DirIn)
    for _,edge := range edges {
            edge.SetProperty("key", []byte("value"))
    }
}

func TestNodesWithKeyValue(t *testing.T) {
	t.Skip()
	//todo
}


// func TestGraphOpenGraph(t *testing.T){
// 	//t.Skip()
// 	graph := core.GetGraph(graphimp)
// 	graph.Open()

// 	if graph == nil {
// 		t.Error("graphdb should not be nil")
// 	} else {
// 		graph.Clear() 
// 		if dbdir != graph.dbdir { t.Error("dbdir not equal")}
// 		if reflect.TypeOf(graph.meta).String() != "*levigo.DB" { t.Error("graph not valid type")}
// 		if graph.meta == nil { t.Error("meta is nil") }
// 		if graph.elements == nil { t.Error("elements is nil") }
// 		if graph.hexaindex == nil { t.Error("hexaindex is nil") }
// 		if graph.props == nil { t.Error("props is nil") }
// 		if bytes.Compare(graph.recsep, []byte("\x1f")) != 0 { t.Error("recsep does not match") }
// 		if graph.EdgeCount() != 0 { t.Error("should have 0 edges")}
// 		if graph.VertexCount() != 0 { t.Errorf("should have 0 vertices has %", graph.VertexCount())}

// 		fi, _ := os.Lstat(dbdir)
// 		if !fi.IsDir() { t.Error("dbdir should be a directory") }
// 		if fi.Name() != "testing.db" { t.Error("dbdir name should match") }
// 		if graph.String() != "#GraphLevigo:dbdir=./testing.db#" { t.Error("String method does not match")}
// 		graph.Close()
// 	}
// }


func TestGraphVertexAdd(t *testing.T) {
	//t.Skip()
	graph := core.GetGraph(graphimp)
	graph.Open()
	graph.Clear()
	defer graph.Close()

	vertex, err := graph.AddVertex(nil)
	assert.True(t, vertex != nil)
	assert.True(t, err == nil)

	id := []byte("somerandomstringid")
	vertex, err = graph.AddVertex(id)
	if assert.True(t, vertex != nil) {
		assert.Equal(t, id, vertex.Id())
		assert.Equal(t, nil, err)
	}
	vertex, err = graph.AddVertex(id) 
	assert.True(t, vertex == nil)
	assert.Equal(t, core.ErrAlreadyExists, err )
}

func TestGraphCloseAndOpen(t *testing.T) {
	//t.Skip()
	graph := core.GetGraph(graphimp)
	dberr := graph.Open()

	if dberr != nil { t.Fatal(dberr) }
	graph.Clear()
	count := graph.EdgeCount()
	if count != 0 { t.Error("should have 0 edges")}
	count = graph.VertexCount()
	if count != 0 { t.Error("should have 0 vertices")}
	graph.AddVertex([]byte("somerandomstringid"))
	count  = graph.VertexCount()
	if count != 1 { t.Error("should have 1 vertex")}
	if graph.Capabilities().Persistent() == true {
		graph.Close()
		graph.Open()
		count  = graph.VertexCount()
		if count != 1 { t.Error("should have 1 vertex")}
	}	
	graph.Close()
}


func TestGraphVertexGet(t *testing.T) {
	//t.Skip()
	graph := core.GetGraph(graphimp)
	graph.Open()
	graph.Clear()
	defer graph.Close()

	ida := []byte("somerandomstringid")
	idb := []byte("idb")
	vertexa, _ := graph.AddVertex(ida)
	vertexb,_  := graph.Vertex(ida)
	assert.Equal(t, vertexa, vertexb)
	testvertex, _ := graph.Vertex(idb)
	assert.True(t, testvertex == nil)
	vertexc, _ := graph.AddVertex(idb)
	vertexd, _ := graph.Vertex(idb)
	assert.Equal(t, vertexc, vertexd)
}


func TestGraphVertexDel(t *testing.T) {
	//t.Skip()
	graph := core.GetGraph(graphimp)
	graph.Open()
	graph.Clear()
	defer graph.Close()

	ida := []byte("somerandomstringid")
	err := graph.DelVertex(nil)
	assert.Equal(t, core.ErrNilValue, err)
	vertex1, _ := graph.AddVertex(ida)
	vertex2, _ := graph.Vertex(ida)
	assert.Equal(t, vertex1, vertex2)
	err = graph.DelVertex(vertex1)
	assert.True(t, err == nil)
	vertex3, _ := graph.Vertex(ida)
	assert.True(t, vertex3 == nil)
	err = graph.DelVertex(vertex1)
	assert.True(t, err == nil)
}

/*

func TestGraphVertexCount(t *testing.T) {
	//t.Skip()
	graph,_ := OpenGraph(dbdir)
	graph.Clear()
	defer graph.Close()

	assert.Equal(t, uint64(0), graph.VertexCount())

	ida := []byte("somerandomstringid")
	vertexa,_ := graph.AddVertex(ida)
	assert.Equal(t, uint64(1), graph.VertexCount())

	testvertii := []*VertexLevigo{}
	var vertex *VertexLevigo
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _ = graph.AddVertex([]byte(a + "-" + n))
			testvertii = append(testvertii, vertex)
		}
	}
	numv := len(testvertii)
	assert.Equal(t, uint64(numv + 1), graph.VertexCount())
	graph.DelVertex(vertexa)
	assert.Equal(t, uint64(numv), graph.VertexCount())
	for i :=0; i < numv; i++ {
		graph.DelVertex(testvertii[i])
		assert.Equal(t, uint64(numv - (i + 1)), graph.VertexCount() )
	}
	assert.Equal(t, uint64(0), graph.VertexCount())
}


func TestGraphVertexGetAll(t *testing.T) {
	//t.Skip()
	graph,_ := OpenGraph(dbdir)
	graph.Clear()
	defer graph.Close()

	ida := []byte("somerandomstringid")
	testvertii := []*VertexLevigo{}
	var vertex *VertexLevigo

	assert.True(t, len(graph.Vertices()) == 0)

	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _= graph.AddVertex([]byte(a + "-" + n))
			testvertii = append(testvertii, vertex)
		}
	}
	assert.Equal(t, testvertii, graph.Vertices())
	lastvertex, _ := graph.AddVertex(ida)
	assert.NotEqual(t, testvertii, graph.Vertices())
	//keys are lexicaly ordered.. lastvertex should be the last in the list
	testvertii = graph.Vertices()
	assert.Equal(t, lastvertex, testvertii[len(testvertii) - 1])
}

func TestGraphEdgeAdd(t *testing.T) {
	//t.Skip()
	graph,_ := OpenGraph(dbdir)
	graph.Clear()
	defer graph.Close()

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")
	vertex1,_ := graph.AddVertex(vid1)
	vertex2,_ := graph.AddVertex(vid2)
	//fmt.Printf("v=%v\n",vertex2)

	edge1, err := graph.AddEdge(nil,vertex1,vertex2,"edgeforward")
	assert.True(t, edge1 == nil)
	assert.Equal(t, NilValue, err)

	edge1, err = graph.AddEdge(eid1,nil,vertex2,"edgeforward")
	assert.True(t, edge1 == nil)
	assert.Equal(t, NilValue, err)

	edge1, err = graph.AddEdge(eid1,vertex1,nil,"edgeforward")
	assert.True(t, edge1 == nil)
	assert.Equal(t, NilValue, err)

	edge1, err = graph.AddEdge(eid1, vertex1, vertex2, "edgeforward")

	if assert.True(t, edge1 != nil) {
		assert.Equal(t, eid1, edge1.Id())
		assert.Equal(t, EdgeType, edge1.Elementtype)
		assert.Equal(t, nil, err)
		assert.Equal(t, "<EdgeLevigo:thisisedge1,s=thisisvertex1,o=thisisvertex2,l=edgeforward@#GraphLevigo:dbdir=./testing.db#>", edge1.String())
		assert.Equal(t, vertex1, edge1.VertexOut())
		assert.Equal(t, vertex2, edge1.VertexIn())
		assert.Equal(t, "edgeforward", edge1.Label())
	}

	edge2, errb := graph.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.True(t, edge2 == nil)
	assert.Equal(t, KeyExists, errb )
}

func TestGraphEdgeGet(t *testing.T) {
	//t.Skip()
	graph,_ := OpenGraph(dbdir)
	graph.Clear()
	defer graph.Close()

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")
	assert.True(t, graph.Edge(eid1) == nil)

	vertex1,_ := graph.AddVertex(vid1)
	vertex2,_ := graph.AddVertex(vid2)
	edge1, _ := graph.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	edge1a := graph.Edge(eid1)
	assert.Equal(t, edge1, edge1a)

	//allow duplicates 
	eid2 := []byte("thisisedge2")
	edge2, _ := graph.AddEdge(eid2, vertex1, vertex2, "edgeforward")
	assert.Equal(t, eid2, edge2.Id())
	assert.Equal(t, vertex1, edge1.VertexOut())
	assert.Equal(t, vertex2, edge1.VertexIn())
	assert.Equal(t, "edgeforward", edge1.Label())
}

func TestGraphEdgeDel(t *testing.T) {
	//t.Skip()
	graph,_ := OpenGraph(dbdir)
	graph.Clear()
	defer graph.Close()

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")

	err := graph.DelEdge(nil)
	assert.Equal(t, NilValue, err)
	err = graph.DelEdge(new(EdgeLevigo))
	assert.Equal(t, NilValue, err)


	edgenull := &EdgeLevigo{new(ElementLevigo),nil,nil,""}
	edgenull.db = graph 
	edgenull.id = eid1 
	edgenull.Elementtype = EdgeType
	err = graph.DelEdge(edgenull)
	assert.Equal(t, KeyDoesNotExist, err)

	vertex1,_ := graph.AddVertex(vid1)
	vertex2,_ := graph.AddVertex(vid2)
	edge1, _ := graph.AddEdge(eid1, vertex1, vertex2, "edgeforward")

	err = graph.DelEdge(edge1)
	assert.True(t, err == nil)
	assert.True(t, graph.Edge(eid1) == nil)
	err = graph.DelEdge(edge1)
	assert.Equal(t, KeyDoesNotExist, err)
}


func TestGraphEdgeGetAll(t *testing.T) {
	//t.Skip()
	graph,_ := OpenGraph(dbdir)
	graph.Clear()
	defer graph.Close()

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")

	vertex1,_ := graph.AddVertex(vid1)
	vertex2,_ := graph.AddVertex(vid2)

	assert.True(t, len(graph.Edges()) == 0)
	testedges := []*EdgeLevigo{}
	var edge *EdgeLevigo
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			edge, _= graph.AddEdge([]byte(a + "-" + n), vertex1, vertex2, "somedge")
			testedges = append(testedges, edge)
		}
	}
	assert.Equal(t, testedges, graph.Edges())
	lastedge, _ := graph.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.NotEqual(t, testedges, graph.Edges())
	//keys are lexicaly ordered.. lastedge should be the last in the list
	testedges = graph.Edges()
	assert.Equal(t, lastedge, testedges[len(testedges) - 1])
}


func TestGraphEdgeCount(t *testing.T) {
	//t.Skip()
	graph,_ := OpenGraph(dbdir)
	graph.Clear()
	defer graph.Close()

	assert.Equal(t, uint64(0), graph.EdgeCount())

	vid1 := []byte("thisisvertex1")
	vid2 := []byte("thisisvertex2")
	eid1 := []byte("thisisedge1")

	vertex1,_ := graph.AddVertex(vid1)
	vertex2,_ := graph.AddVertex(vid2)

	edge1,_ := graph.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.Equal(t, uint64(1), graph.EdgeCount())

	testedges := []*EdgeLevigo{}
	var edge *EdgeLevigo
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			edge, _= graph.AddEdge([]byte(a + "-" + n), vertex1, vertex2, "somedge")
			testedges = append(testedges, edge)
		}
	}
	numv := len(testedges)
	assert.Equal(t, uint64(numv + 1), graph.EdgeCount())
	graph.DelEdge(edge1)
	assert.Equal(t, uint64(numv), graph.EdgeCount())
	for i :=0; i < numv; i++ {
		graph.DelEdge(testedges[i])
		assert.Equal(t, uint64(numv - (i + 1)), graph.EdgeCount() )
	}
	assert.Equal(t, uint64(0), graph.EdgeCount())
}

*/