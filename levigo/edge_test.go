package levigo

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//"os"
	//"reflect"
	//"github.com/jmhodges/levigo"
	//"fmt"
	//"github.com/lexlapax/graveldb/core"	
)


func TestEdge(t *testing.T){
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()
	defer gdb.Close()

	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	vid3 := []byte("vertex3")
	vid4 := []byte("vertex4")
	
	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)
	vertex3,_ := gdb.AddVertex(vid3)
	vertex4,_ := gdb.AddVertex(vid4)

	eid1 := []byte("edge1")
	eid2 := []byte("edge2")
	eid3 := []byte("edge3")
	eid4 := []byte("edge4")

	edge1, _ := gdb.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := gdb.AddEdge(eid2, vertex2, vertex3, "2 to 3")
	edge3, _ := gdb.AddEdge(eid3, vertex3, vertex4, "3 to 4")
	edge4, _ := gdb.AddEdge(eid4, vertex4, vertex1, "4 to 1")

	assert.Equal(t, eid1, edge1.Id())
	assert.Equal(t, vertex1, edge1.VertexOut())
	assert.Equal(t, vertex2, edge1.VertexIn())
	assert.Equal(t, "1 to 2", edge1.Label())

	assert.Equal(t, eid2, edge2.Id())
	assert.Equal(t, vertex2, edge2.VertexOut())
	assert.Equal(t, vertex3, edge2.VertexIn())
	assert.Equal(t, "2 to 3", edge2.Label())


	assert.Equal(t, eid3, edge3.Id())
	assert.Equal(t, vertex3, edge3.VertexOut())
	assert.Equal(t, vertex4, edge3.VertexIn())
	assert.Equal(t, "3 to 4", edge3.Label())


	assert.Equal(t, eid4, edge4.Id())
	assert.Equal(t, vertex4, edge4.VertexOut())
	assert.Equal(t, vertex1, edge4.VertexIn())
	assert.Equal(t, "4 to 1", edge4.Label())

}
