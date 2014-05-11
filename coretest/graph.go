package coretest

import (
	//"strings"
	"github.com/stretchr/testify/assert"
	"github.com/lexlapax/graveldb/core"
)

type GraphTestSuite struct {
	BaseTestSuite
}


func (suite *GraphTestSuite) TestGraphEmpty(){

	//todo - graph capabilites based 
	vertices, _ := suite.TestGraph.Vertices()
	assert.Equal(suite.T(), 0, len(vertices))
	edges, _ := suite.TestGraph.Edges()
	assert.Equal(suite.T(), 0, len(edges))
}

func (suite *GraphTestSuite) TestGraphGuid(){
	uuid := suite.TestGraph.Guid()
	assert.True(suite.T(), uuid != "")
	assert.Equal(suite.T(), 36, len(uuid))
	//assert.Equal(suite.T(), 8, strings.IndexRune(uuid, '-'))
	assert.Equal(suite.T(), "-", string(uuid[8]))
	assert.Equal(suite.T(), "-", string(uuid[13]))
	assert.Equal(suite.T(), "-", string(uuid[18]))
	assert.Equal(suite.T(), "-", string(uuid[23]))
	if suite.TestGraph.Capabilities().Persistent() == true {
		suite.TestGraph.Close()
		suite.TestGraph.Open()
		uuid2  := suite.TestGraph.Guid()
		assert.Equal(suite.T(), uuid, uuid2)
	}	

}


func (suite *GraphTestSuite) TestGraphAdd() {
	a,_ := suite.TestGraph.AddVertex("")
	b,_ := suite.TestGraph.AddVertex("")
    edge, _ := suite.TestGraph.AddEdge("", a, b, "knows")

	//todo - graph capabilites based 
	vertices, _ := suite.TestGraph.Vertices()
	assert.Equal(suite.T(), 2, len(vertices))
	edges, _ := suite.TestGraph.Edges()
	assert.Equal(suite.T(), 1, len(edges))

	suite.TestGraph.DelVertex(a)

	vertices, _ = suite.TestGraph.Vertices()
	assert.Equal(suite.T(), 1, len(vertices))
	edges, _ = suite.TestGraph.Edges()
	assert.Equal(suite.T(), 0, len(edges))
	err := suite.TestGraph.DelEdge(edge)
	assert.True(suite.T(), err != nil)
}

func (suite *GraphTestSuite) TestGraphSetProperties() {
	//todo graph capabilities based
	a,_ := suite.TestGraph.AddVertex("")
	b,_ := suite.TestGraph.AddVertex("")
	suite.TestGraph.AddEdge("", a, b, "knows")
	suite.TestGraph.AddEdge("", a, b, "knows")
	edges, _ := b.Edges(core.DirIn)
    for _,edge := range edges {
            edge.SetProperty("key", []byte("value"))
    }
}


// func (suite *GraphTestSuite) TestGraphOpenGraph(){

// 	if graph == nil {
// 		t.Error("graphdb should not be nil")
// 	} else {
// 		suite.TestGraph.Clear() 
// 		if dbdir != suite.TestGraph.dbdir { t.Error("dbdir not equal")}
// 		if reflect.TypeOf(suite.TestGraph.meta).String() != "*levigo.DB" { t.Error("graph not valid type")}
// 		if suite.TestGraph.meta == nil { t.Error("meta is nil") }
// 		if suite.TestGraph.elements == nil { t.Error("elements is nil") }
// 		if suite.TestGraph.hexaindex == nil { t.Error("hexaindex is nil") }
// 		if suite.TestGraph.props == nil { t.Error("props is nil") }
// 		if bytes.Compare(suite.TestGraph.recsep, []byte("\x1f")) != 0 { t.Error("recsep does not match") }
// 		if suite.TestGraph.EdgeCount() != 0 { t.Error("should have 0 edges")}
// 		if suite.TestGraph.VertexCount() != 0 { t.Errorf("should have 0 vertices has %", suite.TestGraph.VertexCount())}

// 		fi, _ := os.Lstat(dbdir)
// 		if !fi.IsDir() { t.Error("dbdir should be a directory") }
// 		if fi.Name() != "testing.db" { t.Error("dbdir name should match") }
// 		if suite.TestGraph.String() != "#GraphLevigo:dbdir=./testing.db#" { t.Error("String method does not match")}
// 		suite.TestGraph.Close()
// 	}
// }

func (suite *GraphTestSuite) TestGraphVertexAdd() {

	vertex, err := suite.TestGraph.AddVertex("")
	assert.True(suite.T(), vertex != nil)
	assert.True(suite.T(), err == nil)

	id := "somerandomstringid"
	vertex, err = suite.TestGraph.AddVertex(id)
	if assert.True(suite.T(), vertex != nil) {
		assert.Equal(suite.T(), id, vertex.Id())
		assert.Equal(suite.T(), nil, err)
	}
	vertex, err = suite.TestGraph.AddVertex(id) 
	assert.True(suite.T(), vertex == nil)
	assert.Equal(suite.T(), core.ErrAlreadyExists, err )
}

func (suite *GraphTestSuite) TestGraphCloseAndOpen() {

	count := suite.TestGraph.EdgeCount()
	if count != 0 { suite.T().Error("should have 0 edges")}
	count = suite.TestGraph.VertexCount()
	if count != 0 { suite.T().Error("should have 0 vertices")}
	suite.TestGraph.AddVertex("somerandomstringid")
	count  = suite.TestGraph.VertexCount()
	if count != 1 { suite.T().Error("should have 1 vertex")}
	if suite.TestGraph.Capabilities().Persistent() == true {
		suite.TestGraph.Close()
		suite.TestGraph.Open()
		count  = suite.TestGraph.VertexCount()
		if count != 1 { suite.T().Error("should have 1 vertex")}
	}	
}

func (suite *GraphTestSuite) TestGraphVertexGet() {

	ida := "somerandomstringid"
	idb := "idb"
	vertexa, _ := suite.TestGraph.AddVertex(ida)
	vertexb,_  := suite.TestGraph.Vertex(ida)
	assert.Equal(suite.T(), vertexa, vertexb)
	testvertex, _ := suite.TestGraph.Vertex(idb)
	assert.True(suite.T(), testvertex == nil)
	vertexc, _ := suite.TestGraph.AddVertex(idb)
	vertexd, _ := suite.TestGraph.Vertex(idb)
	assert.Equal(suite.T(), vertexc, vertexd)
}

func (suite *GraphTestSuite) TestGraphVertexDel() {

	ida := "somerandomstringid"
	err := suite.TestGraph.DelVertex(nil)
	assert.Equal(suite.T(), core.ErrNilValue, err)
	vertex1, _ := suite.TestGraph.AddVertex(ida)
	vertex2, _ := suite.TestGraph.Vertex(ida)
	assert.Equal(suite.T(), vertex1, vertex2)
	err = suite.TestGraph.DelVertex(vertex1)
	assert.True(suite.T(), err == nil)
	vertex3, _ := suite.TestGraph.Vertex(ida)
	assert.True(suite.T(), vertex3 == nil)
	err = suite.TestGraph.DelVertex(vertex1)
	assert.True(suite.T(), err == nil)
}

func (suite *GraphTestSuite) TestGraphVertexCount() {

	assert.Equal(suite.T(), uint(0), suite.TestGraph.VertexCount())

	ida := "somerandomstringid"
	vertexa,_ := suite.TestGraph.AddVertex(ida)
	assert.Equal(suite.T(), uint(1), suite.TestGraph.VertexCount())

	testvertii := []core.Vertex{}
	var vertex core.Vertex
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _ = suite.TestGraph.AddVertex(a + "-" + n)
			testvertii = append(testvertii, vertex)
		}
	}
	numv := len(testvertii)
	assert.Equal(suite.T(), uint(numv + 1), suite.TestGraph.VertexCount())
	suite.TestGraph.DelVertex(vertexa)
	assert.Equal(suite.T(), uint(numv), suite.TestGraph.VertexCount())
	for i :=0; i < numv; i++ {
		suite.TestGraph.DelVertex(testvertii[i])
		assert.Equal(suite.T(), uint(numv - (i + 1)), suite.TestGraph.VertexCount() )
	}
	assert.Equal(suite.T(), uint(0), suite.TestGraph.VertexCount())
}

func (suite *GraphTestSuite) TestGraphVertexGetAll() {

	ida := "somerandomstringid"
	testvertii := []core.Vertex{}
	var vertex core.Vertex

	assert.True(suite.T(), suite.TestGraph.VertexCount() == uint(0))

	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _= suite.TestGraph.AddVertex(a + "-" + n)
			testvertii = append(testvertii, vertex)
		}
	}
	verticesget,_  := suite.TestGraph.Vertices()
	assert.Equal(suite.T(), len(testvertii), len(verticesget))
	// assert.Equal(suite.T(), testvertii, verticesget)
	suite.TestGraph.AddVertex(ida)
	verticesget,_  = suite.TestGraph.Vertices()
	// assert.NotEqual(suite.T(), testvertii, verticesget)

	assert.Equal(suite.T(), len(testvertii) + 1, len(verticesget))

	// //keys are lexicaly ordered.. lastvertex should be the last in the list
	// verticesget = suite.TestGraph.Vertices()
	// assert.Equal(suite.T(), lastvertex, verticesget[len(verticesget) - 1])
}

func (suite *GraphTestSuite) TestGraphIterVertices() {

	atomset := core.NewAtomSet()
	atomset.Clear()
	var vertex core.Vertex

	assert.True(suite.T(), suite.TestGraph.VertexCount() == uint(0))

	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _= suite.TestGraph.AddVertex(a + "-" + n)
			atomset.Add(vertex)
		}
	}
	testatomset := core.NewAtomSet()
	testatomset.Clear()
	for vertex = range suite.TestGraph.IterVertices() {
		testatomset.Add(vertex)
		assert.True(suite.T(), atomset.Contains(vertex))
	} 	

	assert.Equal(suite.T(), atomset.Count(), testatomset.Count())
}


func (suite *GraphTestSuite) TestGraphEdgeAdd() {

	vid1 := "thisisvertex1"
	vid2 := "thisisvertex2"
	eid1 := "thisisedge1"
	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)
	//fmt.Printf("v=%v\n",vertex2)

	// edge1, err := suite.TestGraph.AddEdge(nil,vertex1,vertex2,"edgeforward")
	// assert.True(suite.T(), edge1 != nil)

	edge1, err := suite.TestGraph.AddEdge(eid1,nil,vertex2,"edgeforward")
	assert.True(suite.T(), edge1 == nil)
	assert.Equal(suite.T(), core.ErrNilValue, err)

	edge1, err = suite.TestGraph.AddEdge(eid1,vertex1,nil,"edgeforward")
	assert.True(suite.T(), edge1 == nil)
	assert.Equal(suite.T(), core.ErrNilValue, err)

	edge1, err = suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "edgeforward")

	if assert.True(suite.T(), edge1 != nil) {
		assert.Equal(suite.T(), eid1, edge1.Id())
		assert.Equal(suite.T(), nil, err)
		testvertex,_ :=  edge1.VertexOut()
		assert.Equal(suite.T(), vertex1, testvertex)
		testvertex,_ =  edge1.VertexIn()
		assert.Equal(suite.T(), vertex2, testvertex)
		assert.Equal(suite.T(), "edgeforward", edge1.Label())
	}

	edge2, errb := suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.True(suite.T(), edge2 == nil)
	assert.Equal(suite.T(), core.ErrAlreadyExists, errb )
}

func (suite *GraphTestSuite) TestGraphEdgeGet() {
	vid1 := "thisisvertex1"
	vid2 := "thisisvertex2"
	eid1 := "thisisedge1"
	edge1, _ := suite.TestGraph.Edge(eid1)
	assert.True(suite.T(), edge1 == nil)

	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)
	edge1, _ = suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	edge1a,_ := suite.TestGraph.Edge(eid1)
	assert.Equal(suite.T(), edge1, edge1a)

	//allow duplicates 
	eid2 := "thisisedge2"
	edge2, _ := suite.TestGraph.AddEdge(eid2, vertex1, vertex2, "edgeforward")
	assert.Equal(suite.T(), eid2, edge2.Id())
	testvertex,_ := edge1.VertexOut()
	assert.Equal(suite.T(), vertex1, testvertex)
	testvertex,_ = edge1.VertexIn()
	assert.Equal(suite.T(), vertex2, testvertex)
	assert.Equal(suite.T(), "edgeforward", edge1.Label())
	assert.Equal(suite.T(), "edgeforward", edge2.Label())
}


func (suite *GraphTestSuite) TestGraphEdgeDel() {
	vid1 := "thisisvertex1"
	vid2 := "thisisvertex2"
	eid1 := "thisisedge1"

	err := suite.TestGraph.DelEdge(nil)
	assert.Equal(suite.T(), core.ErrNilValue, err)

	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)
	edge1, _ := suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "edgeforward")

	err = suite.TestGraph.DelEdge(edge1)
	assert.True(suite.T(), err == nil)
	testedge,_ := suite.TestGraph.Edge(eid1)
	assert.True(suite.T(),  testedge == nil)
	err = suite.TestGraph.DelEdge(edge1)
	assert.Equal(suite.T(), core.ErrDoesntExist, err)
}

func (suite *GraphTestSuite) TestGraphEdgeCount() {
	assert.Equal(suite.T(), uint(0), suite.TestGraph.EdgeCount())

	vid1 := "thisisvertex1"
	vid2 := "thisisvertex2"
	eid1 := "thisisedge1"

	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)

	edge1,_ := suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	assert.Equal(suite.T(), uint(1), suite.TestGraph.EdgeCount())

	testedges := []core.Edge{}
	var edge core.Edge
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			edge, _= suite.TestGraph.AddEdge(a + "-" + n, vertex1, vertex2, "somedge")
			testedges = append(testedges, edge)
		}
	}
	numv := len(testedges)
	assert.Equal(suite.T(), uint64(numv + 1), suite.TestGraph.EdgeCount())
	suite.TestGraph.DelEdge(edge1)
	assert.Equal(suite.T(), uint64(numv), suite.TestGraph.EdgeCount())
	for i :=0; i < numv; i++ {
		suite.TestGraph.DelEdge(testedges[i])
		assert.Equal(suite.T(), uint(numv - (i + 1)), suite.TestGraph.EdgeCount() )
	}
	assert.Equal(suite.T(), uint64(0), suite.TestGraph.EdgeCount())
}

func (suite *GraphTestSuite) TestGraphEdgeGetAll() {

	vid1 := "thisisvertex1"
	vid2 := "thisisvertex2"
	eid1 := "thisisedge1"

	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)

	assert.True(suite.T(), suite.TestGraph.EdgeCount() == uint(0))
	// testedges := core.NewEdgeSet()
	testedges := core.NewAtomSet()
	var edge core.Edge
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			edge, _= suite.TestGraph.AddEdge(a + "-" + n, vertex1, vertex2, "somedge")
			testedges.Add(edge)
		}
	}
	edges, _ := suite.TestGraph.Edges()
	assert.Equal(suite.T(), testedges.Count(), len(edges))
	for _,e := range edges {
		assert.True(suite.T(), testedges.Contains(e))
	}
	suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "edgeforward")
	edges, _ = suite.TestGraph.Edges()
	assert.Equal(suite.T(), testedges.Count() + 1, len(edges))
}

func (suite *GraphTestSuite) TestGraphIterEdges() {

	vid1 := "thisisvertex1"
	vid2 := "thisisvertex2"

	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)

	assert.True(suite.T(), suite.TestGraph.EdgeCount() == uint(0))
	// testedges := core.NewEdgeSet()
	testedges := core.NewAtomSet()
	var edge core.Edge
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			edge, _= suite.TestGraph.AddEdge(a + "-" + n, vertex1, vertex2, "somedge")
			testedges.Add(edge)
		}
	}

	edges := core.NewAtomSet()
	edges.Clear()
	for edge = range suite.TestGraph.IterEdges() {
		edges.Add(edge)
		assert.True(suite.T(), testedges.Contains(edge))
	} 	

	assert.Equal(suite.T(), testedges.Count(), edges.Count())
}
