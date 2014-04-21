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

func TestElementProperty(t *testing.T){
	//t.Skip()
	gdb,_ := OpenGraph(dbdir)
	gdb.Clear()

	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	
	vertex1,_ := gdb.AddVertex(vid1)
	vertex2,_ := gdb.AddVertex(vid2)

	eid1 := []byte("edge1")
	eid2 := []byte("edge2")

	edge1, _ := gdb.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := gdb.AddEdge(eid2, vertex2, vertex1, "2 to 1")

	assert.True(t, vertex1.Property("name") == nil)
	assert.True(t,  vertex2.Property("name") == nil)
	assert.True(t, edge1.Property("name") == nil)
	assert.True(t, edge2.Property("name") == nil)

	vertex1.SetProperty("name", nil)
	assert.Equal(t, []byte{}, vertex1.Property("name"))	
	vertex1.SetProperty("name", []byte("some value"))
	assert.Equal(t, []byte("some value"), vertex1.Property("name"))


	vertex1.SetProperty("name", []byte("this is vertex 1"))
	vertex1.SetProperty("nickname", []byte("node 1"))
	vertex2.SetProperty("name", []byte("this is vertex 2"))
	vertex2.SetProperty("nickname", []byte("node 2"))
	edge1.SetProperty("name", []byte("this is edge 1"))
	edge1.SetProperty("nickname", []byte("connector 1"))
	edge2.SetProperty("name", []byte("this is edge 2"))
	edge2.SetProperty("nickname", []byte("connector 2"))


	assert.Equal(t, []byte("this is vertex 1"), vertex1.Property("name"))
	assert.Equal(t, []byte("node 1"), vertex1.Property("nickname"))
	assert.Equal(t, []byte("this is vertex 2"), vertex2.Property("name"))
	assert.Equal(t, []byte("node 2"), vertex2.Property("nickname"))
	assert.Equal(t, []byte("this is edge 1"), edge1.Property("name"))
	assert.Equal(t, []byte("connector 1"), edge1.Property("nickname"))
	assert.Equal(t, []byte("this is edge 2"), edge2.Property("name"))
	assert.Equal(t, []byte("connector 2"), edge2.Property("nickname"))

	assert.True(t, vertex1.DelProperty("nonexistentproperty") == nil)
	assert.Equal(t, []byte("this is vertex 1"), vertex1.DelProperty("name"))
	assert.True(t, vertex1.Property("name") == nil)
	assert.Equal(t, []byte("node 1"), vertex1.DelProperty("nickname"))
	assert.True(t, vertex1.Property("nickname") == nil)

	gdb.Close()

	gdb,_ = OpenGraph(dbdir)
	vertex1 = gdb.Vertex(vid1)
	vertex2 = gdb.Vertex(vid2)
	edge1 = gdb.Edge(eid1)
	edge2 = gdb.Edge(eid2)
	assert.True(t, vertex1.Property("name") == nil)
	assert.True(t, vertex1.Property("nickname") == nil)
	assert.Equal(t, []byte("this is vertex 2"), vertex2.Property("name"))
	assert.Equal(t, []byte("node 2"), vertex2.Property("nickname"))
	assert.Equal(t, []byte("this is edge 1"), edge1.Property("name"))
	assert.Equal(t, []byte("connector 1"), edge1.Property("nickname"))
	assert.Equal(t, []byte("this is edge 2"), edge2.Property("name"))
	assert.Equal(t, []byte("connector 2"), edge2.Property("nickname"))

	propkeys := []string{}
	var propkey string
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			propkey = a + "-" + n
			propkeys = append(propkeys, propkey)
			vertex1.SetProperty(propkey, []byte(propkey + "-vertex"))
			edge1.SetProperty(propkey, []byte(propkey + "-edge"))
		}
	}
	assert.Equal(t, propkeys, vertex1.PropertyKeys())
	for _, key := range propkeys {
		propvertexval := key + "-vertex"
		propedgeval := key + "-edge"
		assert.Equal(t, []byte(propvertexval), vertex1.Property(key))
		assert.Equal(t, []byte(propedgeval), edge1.Property(key))
	}

	propkeys = append(propkeys, "name")
	propkeys = append(propkeys, "nickname")
	assert.Equal(t, propkeys, edge1.PropertyKeys())

	gdb.Close()

}
