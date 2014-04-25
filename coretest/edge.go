package coretest

import (
	"github.com/stretchr/testify/assert"
	"github.com/lexlapax/graveldb/core"
)

type EdgeTestSuite struct {
	BaseTestSuite
}

func (suite *EdgeTestSuite) TestEdge(){
	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	vid3 := []byte("vertex3")
	vid4 := []byte("vertex4")
	
	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)
	vertex3,_ := suite.TestGraph.AddVertex(vid3)
	vertex4,_ := suite.TestGraph.AddVertex(vid4)

	eid1 := []byte("edge1")
	eid2 := []byte("edge2")
	eid3 := []byte("edge3")
	eid4 := []byte("edge4")

	edge1, _ := suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := suite.TestGraph.AddEdge(eid2, vertex2, vertex3, "2 to 3")
	edge3, _ := suite.TestGraph.AddEdge(eid3, vertex3, vertex4, "3 to 4")
	edge4, _ := suite.TestGraph.AddEdge(eid4, vertex4, vertex1, "4 to 1")

	assert.Equal(suite.T(), eid1, edge1.Id())
	testvertex,_ := edge1.VertexOut()
	assert.Equal(suite.T(), vertex1, testvertex)
	testvertex,_ = edge1.VertexIn()
	assert.Equal(suite.T(), vertex2, testvertex)
	assert.Equal(suite.T(), "1 to 2", edge1.Label())

	assert.Equal(suite.T(), eid2, edge2.Id())
	testvertex,_ = edge2.Vertex(core.DirOut)
	assert.Equal(suite.T(), vertex2, testvertex)
	testvertex,_ = edge2.Vertex(core.DirIn)
	assert.Equal(suite.T(), vertex3, testvertex)
	assert.Equal(suite.T(), "2 to 3", edge2.Label())


	assert.Equal(suite.T(), eid3, edge3.Id())
	testvertex,_ = edge3.VertexOut()
	assert.Equal(suite.T(), vertex3, testvertex)
	testvertex,_ = edge3.VertexIn()
	assert.Equal(suite.T(), vertex4, testvertex)
	assert.Equal(suite.T(), "3 to 4", edge3.Label())


	assert.Equal(suite.T(), eid4, edge4.Id())
	testvertex,_ = edge4.VertexOut()
	assert.Equal(suite.T(), vertex4, testvertex)
	testvertex,_ = edge4.VertexIn()
	assert.Equal(suite.T(), vertex1, testvertex)
	assert.Equal(suite.T(), "4 to 1", edge4.Label())
	_, err := edge4.Vertex(core.DirAny)
	assert.True(suite.T(), err != nil)

}
