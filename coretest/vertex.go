package coretest

import (
	"github.com/stretchr/testify/assert"
	"github.com/lexlapax/graveldb/core"
)

type VertexTestSuite struct {
	BaseTestSuite
}

func (suite *VertexTestSuite) TestVertex(){

	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	vid3 := []byte("vertex3")
	vid4 := []byte("vertex4")

	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)
	vertex3,_ := suite.TestGraph.AddVertex(vid3)
	vertex4,_ := suite.TestGraph.AddVertex(vid4)

	// this is what we will test
	// v1 points to v2 and v3 which both point to v4 which points back to v1

	// 		v2
	// 	/		\
	// v1			v4 - v1
	// 	\		/
	// 		v3
	
	//

	edges, _ := vertex1.OutEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges, _ = vertex2.OutEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges, _ = vertex3.OutEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges, _ = vertex4.OutEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges, _ = vertex1.InEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges, _ = vertex2.InEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges, _ = vertex3.InEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges, _ = vertex4.InEdges()
	assert.True(suite.T(), len(edges) == 0)

	eid1 := []byte("edge1")
	eid2 := []byte("edge2")
	eid3 := []byte("edge3")
	eid4 := []byte("edge4")
	eid5 := []byte("edge5")

	edge1, _ := suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := suite.TestGraph.AddEdge(eid2, vertex1, vertex3, "1 to 3")
	edge3, _ := suite.TestGraph.AddEdge(eid3, vertex2, vertex4, "2 to 4")
	edge4, _ := suite.TestGraph.AddEdge(eid4, vertex3, vertex4, "3 to 4")
	edge5, _ := suite.TestGraph.AddEdge(eid5, vertex4, vertex1, "4 to 1")

	edges, _ = suite.TestGraph.Edges()
	assert.True(suite.T(), len(edges) == 5)

	edges, _ = vertex1.Edges(core.DirAny)
	assert.True(suite.T(), len(edges) == 3)

	edges, _ = vertex4.Edges(core.DirAny)
	assert.True(suite.T(), len(edges) == 3)

	v1out, _ := vertex1.OutEdges()
	v1in, _  := vertex1.InEdges()
	v2out, _ := vertex2.OutEdges()
	v2in, _  := vertex2.InEdges()
	v3out, _ := vertex3.OutEdges()
	v3in, _  := vertex3.InEdges()
	v4out, _ := vertex4.OutEdges()
	v4in, _  := vertex4.InEdges()

	assert.True(suite.T(), 2 == len(v1out))
	assert.True(suite.T(), 1 == len(v1in))
	assert.True(suite.T(), 1 == len(v2out))
	assert.True(suite.T(), 1 == len(v2in))
	assert.True(suite.T(), 1 == len(v3out))
	assert.True(suite.T(), 1 == len(v3in))
	assert.True(suite.T(), 1 == len(v4out))
	assert.True(suite.T(), 2 == len(v4in))

	assert.Equal(suite.T(), edge1, v1out[0])
	assert.Equal(suite.T(), edge2, v1out[1])
	assert.Equal(suite.T(), edge5, v1in[0])

	assert.Equal(suite.T(), edge3, v2out[0])
	assert.Equal(suite.T(), edge1, v2in[0])

	assert.Equal(suite.T(), edge4, v3out[0])
	assert.Equal(suite.T(), edge2, v3in[0])

	assert.Equal(suite.T(), edge5, v4out[0])
	assert.Equal(suite.T(), edge3, v4in[0])
	assert.Equal(suite.T(), edge4, v4in[1])

	err := suite.TestGraph.DelEdge(edge2)
	assert.True(suite.T(), err == nil)
	testedge,_ := suite.TestGraph.Edge(eid2)
	assert.True(suite.T(), testedge == nil)

	edges,_ = suite.TestGraph.Edges()
	assert.True(suite.T(), len(edges) == 4)

	edges,_ = vertex1.Edges(core.DirOut)
	//fmt.Printf("v1out=%v\n", edges)
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex2.Edges(core.DirOut)
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex3.Edges(core.DirOut)
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex4.Edges(core.DirOut)
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex1.Edges(core.DirIn)
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex2.Edges(core.DirIn)
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex3.Edges(core.DirIn)
	assert.True(suite.T(), len(edges) == 0)
	edges,_ = vertex4.Edges(core.DirIn)
	assert.True(suite.T(), len(edges) == 2)


	suite.TestGraph.DelVertex(vertex4)
	edges,_ = suite.TestGraph.Edges()
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex1.OutEdges()
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex2.OutEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges,_ = vertex3.OutEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges,_ = vertex1.InEdges()
	assert.True(suite.T(), len(edges) == 0)
	edges,_ = vertex2.InEdges()
	assert.True(suite.T(), len(edges) == 1)
	edges,_ = vertex3.InEdges()
	assert.True(suite.T(), len(edges) == 0)

	if suite.TestGraph.Capabilities().Persistent() == true {
		suite.TestGraph.Close()	

		suite.TestGraph.Open()
		edges,_ = suite.TestGraph.Edges()
		vertices,_ := suite.TestGraph.Vertices()
		assert.True(suite.T(), len(edges) == 1)
		assert.True(suite.T(), len(vertices) == 3)

		vertex1,_ = suite.TestGraph.Vertex(vid1)
		vertex2,_ = suite.TestGraph.Vertex(vid2)
		vertex3,_ = suite.TestGraph.Vertex(vid3)
		edge1,_ = suite.TestGraph.Edge(eid1)
		edge2,_ = suite.TestGraph.Edge(eid2)
		edges,_ = vertex1.OutEdges()
		assert.True(suite.T(), len(edges) == 1)
		edges,_ = vertex2.OutEdges()
		assert.True(suite.T(), len(edges) == 0)
		edges,_ = vertex3.OutEdges()
		assert.True(suite.T(), len(edges) == 0)
		edges,_ = vertex1.InEdges()
		assert.True(suite.T(), len(edges) == 0)
		edges,_ = vertex2.InEdges()
		assert.True(suite.T(), len(edges) == 1)
		edges,_ = vertex3.InEdges()
		assert.True(suite.T(), len(edges) == 0)

		assert.True(suite.T(), edge2 == nil)

	}
	assert.Equal(suite.T(), "1 to 2", edge1.Label())
	edges, _ = vertex1.OutEdges()
	assert.Equal(suite.T(), edge1, edges[0])
	edges, _ = vertex2.InEdges()
	assert.Equal(suite.T(), edge1, edges[0])
}
