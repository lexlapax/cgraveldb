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


func TestVertex(t *testing.T){
	//t.Skip()
	graph := core.GetGraph(graphimp)
	graph.Open()
	graph.Clear()

	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	vid3 := []byte("vertex3")
	vid4 := []byte("vertex4")

	vertex1,_ := graph.AddVertex(vid1)
	vertex2,_ := graph.AddVertex(vid2)
	vertex3,_ := graph.AddVertex(vid3)
	vertex4,_ := graph.AddVertex(vid4)

	// this is what we will test
	// v1 points to v2 and v3 which both point to v4 which points back to v1

	// 		v2
	// 	/		\
	// v1			v4 - v1
	// 	\		/
	// 		v3
	
	//

	edges, _ := vertex1.OutEdges()
	assert.True(t, len(edges) == 0)
	edges, _ = vertex2.OutEdges()
	assert.True(t, len(edges) == 0)
	edges, _ = vertex3.OutEdges()
	assert.True(t, len(edges) == 0)
	edges, _ = vertex4.OutEdges()
	assert.True(t, len(edges) == 0)
	edges, _ = vertex1.InEdges()
	assert.True(t, len(edges) == 0)
	edges, _ = vertex2.InEdges()
	assert.True(t, len(edges) == 0)
	edges, _ = vertex3.InEdges()
	assert.True(t, len(edges) == 0)
	edges, _ = vertex4.InEdges()
	assert.True(t, len(edges) == 0)

	eid1 := []byte("edge1")
	eid2 := []byte("edge2")
	eid3 := []byte("edge3")
	eid4 := []byte("edge4")
	eid5 := []byte("edge5")

	edge1, _ := graph.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := graph.AddEdge(eid2, vertex1, vertex3, "1 to 3")
	edge3, _ := graph.AddEdge(eid3, vertex2, vertex4, "2 to 4")
	edge4, _ := graph.AddEdge(eid4, vertex3, vertex4, "3 to 4")
	edge5, _ := graph.AddEdge(eid5, vertex4, vertex1, "4 to 1")

	edges, _ = graph.Edges()
	assert.True(t, len(edges) == 5)

	edges, _ = vertex1.Edges(core.DirAny)
	assert.True(t, len(edges) == 3)

	edges, _ = vertex4.Edges(core.DirAny)
	assert.True(t, len(edges) == 3)

	v1out, _ := vertex1.OutEdges()
	v1in, _  := vertex1.InEdges()
	v2out, _ := vertex2.OutEdges()
	v2in, _  := vertex2.InEdges()
	v3out, _ := vertex3.OutEdges()
	v3in, _  := vertex3.InEdges()
	v4out, _ := vertex4.OutEdges()
	v4in, _  := vertex4.InEdges()

	assert.True(t, 2 == len(v1out))
	assert.True(t, 1 == len(v1in))
	assert.True(t, 1 == len(v2out))
	assert.True(t, 1 == len(v2in))
	assert.True(t, 1 == len(v3out))
	assert.True(t, 1 == len(v3in))
	assert.True(t, 1 == len(v4out))
	assert.True(t, 2 == len(v4in))

	assert.Equal(t, edge1, v1out[0])
	assert.Equal(t, edge2, v1out[1])
	assert.Equal(t, edge5, v1in[0])

	assert.Equal(t, edge3, v2out[0])
	assert.Equal(t, edge1, v2in[0])

	assert.Equal(t, edge4, v3out[0])
	assert.Equal(t, edge2, v3in[0])

	assert.Equal(t, edge5, v4out[0])
	assert.Equal(t, edge3, v4in[0])
	assert.Equal(t, edge4, v4in[1])

	err := graph.DelEdge(edge2)
	assert.True(t, err == nil)
	testedge,_ := graph.Edge(eid2)
	assert.True(t, testedge == nil)

	edges,_ = graph.Edges()
	assert.True(t, len(edges) == 4)

	edges,_ = vertex1.Edges(core.DirOut)
	//fmt.Printf("v1out=%v\n", edges)
	assert.True(t, len(edges) == 1)
	edges,_ = vertex2.Edges(core.DirOut)
	assert.True(t, len(edges) == 1)
	edges,_ = vertex3.Edges(core.DirOut)
	assert.True(t, len(edges) == 1)
	edges,_ = vertex4.Edges(core.DirOut)
	assert.True(t, len(edges) == 1)
	edges,_ = vertex1.Edges(core.DirIn)
	assert.True(t, len(edges) == 1)
	edges,_ = vertex2.Edges(core.DirIn)
	assert.True(t, len(edges) == 1)
	edges,_ = vertex3.Edges(core.DirIn)
	assert.True(t, len(edges) == 0)
	edges,_ = vertex4.Edges(core.DirIn)
	assert.True(t, len(edges) == 2)


	graph.DelVertex(vertex4)
	edges,_ = graph.Edges()
	assert.True(t, len(edges) == 1)
	edges,_ = vertex1.OutEdges()
	assert.True(t, len(edges) == 1)
	edges,_ = vertex2.OutEdges()
	assert.True(t, len(edges) == 0)
	edges,_ = vertex3.OutEdges()
	assert.True(t, len(edges) == 0)
	edges,_ = vertex1.InEdges()
	assert.True(t, len(edges) == 0)
	edges,_ = vertex2.InEdges()
	assert.True(t, len(edges) == 1)
	edges,_ = vertex3.InEdges()
	assert.True(t, len(edges) == 0)

	if graph.Capabilities().Persistent() == true {
		graph.Close()	

		graph.Open()
		edges,_ = graph.Edges()
		vertices,_ := graph.Vertices()
		assert.True(t, len(edges) == 1)
		assert.True(t, len(vertices) == 3)

		vertex1,_ = graph.Vertex(vid1)
		vertex2,_ = graph.Vertex(vid2)
		vertex3,_ = graph.Vertex(vid3)
		edge1,_ = graph.Edge(eid1)
		edge2,_ = graph.Edge(eid2)
		edges,_ = vertex1.OutEdges()
		assert.True(t, len(edges) == 1)
		edges,_ = vertex2.OutEdges()
		assert.True(t, len(edges) == 0)
		edges,_ = vertex3.OutEdges()
		assert.True(t, len(edges) == 0)
		edges,_ = vertex1.InEdges()
		assert.True(t, len(edges) == 0)
		edges,_ = vertex2.InEdges()
		assert.True(t, len(edges) == 1)
		edges,_ = vertex3.InEdges()
		assert.True(t, len(edges) == 0)

		assert.True(t, edge2 == nil)

	}
	assert.Equal(t, "1 to 2", edge1.Label())
	edges, _ = vertex1.OutEdges()
	assert.Equal(t, edge1, edges[0])
	edges, _ = vertex2.InEdges()
	assert.Equal(t, edge1, edges[0])
	graph.Close()	
}
