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
func init() {
	mem.Register()
}

func TestEdge(t *testing.T){
	//t.Skip()
	graph := core.GetGraph(graphimp)
	graph.Open()
	defer graph.Close()

	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	vid3 := []byte("vertex3")
	vid4 := []byte("vertex4")
	
	vertex1,_ := graph.AddVertex(vid1)
	vertex2,_ := graph.AddVertex(vid2)
	vertex3,_ := graph.AddVertex(vid3)
	vertex4,_ := graph.AddVertex(vid4)

	eid1 := []byte("edge1")
	eid2 := []byte("edge2")
	eid3 := []byte("edge3")
	eid4 := []byte("edge4")

	edge1, _ := graph.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := graph.AddEdge(eid2, vertex2, vertex3, "2 to 3")
	edge3, _ := graph.AddEdge(eid3, vertex3, vertex4, "3 to 4")
	edge4, _ := graph.AddEdge(eid4, vertex4, vertex1, "4 to 1")

	assert.Equal(t, eid1, edge1.Id())
	testvertex,_ := edge1.VertexOut()
	assert.Equal(t, vertex1, testvertex)
	testvertex,_ = edge1.VertexIn()
	assert.Equal(t, vertex2, testvertex)
	assert.Equal(t, "1 to 2", edge1.Label())

	assert.Equal(t, eid2, edge2.Id())
	testvertex,_ = edge2.Vertex(core.DirOut)
	assert.Equal(t, vertex2, testvertex)
	testvertex,_ = edge2.Vertex(core.DirIn)
	assert.Equal(t, vertex3, testvertex)
	assert.Equal(t, "2 to 3", edge2.Label())


	assert.Equal(t, eid3, edge3.Id())
	testvertex,_ = edge3.VertexOut()
	assert.Equal(t, vertex3, testvertex)
	testvertex,_ = edge3.VertexIn()
	assert.Equal(t, vertex4, testvertex)
	assert.Equal(t, "3 to 4", edge3.Label())


	assert.Equal(t, eid4, edge4.Id())
	testvertex,_ = edge4.VertexOut()
	assert.Equal(t, vertex4, testvertex)
	testvertex,_ = edge4.VertexIn()
	assert.Equal(t, vertex1, testvertex)
	assert.Equal(t, "4 to 1", edge4.Label())
	_, err := edge4.Vertex(core.DirAny)
	assert.True(t, err != nil)

}
