package levelgraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//"os"
	//"reflect"
	//"github.com/jmhodges/levigo"
	//"fmt"
	//"github.com/lexlapax/graveldb/core"	
)

func TestVertex(t *testing.T){
	t.Skip()
	cleanup(dbdir)
	defer cleanup(dbdir)
	gdb,_ := opengraph(dbdir)

	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	vid3 := []byte("vertex3")
	vid4 := []byte("vertex4")

	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)
	vertex3,_ := gdb.AddVertex(vid3)
	vertex4,_ := gdb.AddVertex(vid4)

	/*
	this is what we will test
	v1 points to v2 and v3 which both point to v4 which points back to v1

			v2
		/		\
	v1			v4 - v1
		\		/
			v3
	
	*/


	assert.True(t, vertex1.OutEdges() == nil)
	assert.True(t, vertex2.OutEdges() == nil)
	assert.True(t, vertex3.OutEdges() == nil)
	assert.True(t, vertex4.OutEdges() == nil)
	assert.True(t, vertex1.InEdges() == nil)
	assert.True(t, vertex2.InEdges() == nil)
	assert.True(t, vertex3.InEdges() == nil)
	assert.True(t, vertex4.InEdges() == nil)

	
	eid1 := []byte("edge1")
	eid2 := []byte("edge2")
	eid3 := []byte("edge3")
	eid4 := []byte("edge4")
	eid5 := []byte("edge5")

	edge1, _ := gdb.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := gdb.AddEdge(eid2, vertex1, vertex3, "1 to 3")
	edge3, _ := gdb.AddEdge(eid3, vertex2, vertex4, "2 to 4")
	edge4, _ := gdb.AddEdge(eid4, vertex3, vertex4, "3 to 4")
	edge5, _ := gdb.AddEdge(eid5, vertex4, vertex1, "4 to 1")

	assert.True(t, len(gdb.Edges()) == 5)

	v1out := vertex1.OutEdges()
	v1in  := vertex1.InEdges()
	v2out := vertex2.OutEdges()
	v2in  := vertex2.InEdges()
	v3out := vertex3.OutEdges()
	v3in  := vertex3.InEdges()
	v4out := vertex4.OutEdges()
	v4in  := vertex4.InEdges()

	assert.True(t, 2 == len(v1out))
	assert.True(t, 1 == len(v1in))
	assert.True(t, 1 == len(v2out))
	assert.True(t, 1 == len(v2in))
	assert.True(t, 1 == len(v3out))
	assert.True(t, 2 == len(v3in))
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

	gdb.DelVertex(vertex4)
	assert.True(t, len(gdb.Edges()) == 2)
	assert.True(t, len(vertex1.OutEdges()) == 2)
	assert.True(t, vertex2.OutEdges() == nil)
	assert.True(t, vertex3.OutEdges() == nil)
	assert.True(t, vertex1.InEdges() == nil)
	assert.True(t, len(vertex2.InEdges()) == 1)
	assert.True(t, len(vertex3.InEdges()) == 1)
	gdb.Close()	

	gdb,_ = opengraph(dbdir)
	vertex1 = gdb.Vertex(vid1)
	vertex1 = gdb.Vertex(vid1)
	vertex1 = gdb.Vertex(vid1)
	edge1 = gdb.Edge(eid1)
	edge2 = gdb.Edge(eid2)
	assert.True(t, len(gdb.Edges()) == 2)
	assert.True(t, len(vertex1.OutEdges()) == 2)
	assert.True(t, vertex2.OutEdges() == nil)
	assert.True(t, vertex3.OutEdges() == nil)
	assert.True(t, vertex1.InEdges() == nil)
	assert.True(t, len(vertex2.InEdges()) == 1)
	assert.True(t, len(vertex3.InEdges()) == 1)
	assert.Equal(t, "1 to 2", edge1.Label())
	assert.Equal(t, "1 to 3", edge2.Label())

	gdb.DelEdge(edge2)
	assert.True(t, len(gdb.Edges()) == 1)
	assert.True(t, vertex3.InEdges() == nil)
	assert.True(t, len(vertex1.OutEdges()) == 1)
	assert.True(t, len(vertex2.InEdges()) == 1)
	assert.Equal(t, edge1, vertex1.OutEdges()[0])
	assert.Equal(t, edge1, vertex2.InEdges()[0])

	gdb.Close()	
}
