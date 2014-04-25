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

func TestAtomProperty(t *testing.T){
	//t.Skip()
	graph := core.GetGraph(graphimp)
	graph.Open()

	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	
	vertex1,_ := graph.AddVertex(vid1)
	vertex2,_ := graph.AddVertex(vid2)

	eid1 := []byte("edge1")
	eid2 := []byte("edge2")

	edge1, _ := graph.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := graph.AddEdge(eid2, vertex2, vertex1, "2 to 1")

	propval,_ := vertex1.Property("name")
	assert.True(t,  propval == nil)
	propval,_ = vertex2.Property("name")
	assert.True(t,  propval == nil)
	propval,_ = edge1.Property("name")
	assert.True(t, propval == nil)
	propval,_ = edge2.Property("name")
	assert.True(t, propval == nil)

	vertex1.SetProperty("name", nil)
	propval,_ = vertex1.Property("name")
	assert.Equal(t, []byte{}, propval)	
	vertex1.SetProperty("name", []byte("some value"))
	propval,_ = vertex1.Property("name")
	assert.Equal(t, []byte("some value"), propval)


	vertex1.SetProperty("name", []byte("this is vertex 1"))
	vertex1.SetProperty("nickname", []byte("node 1"))
	vertex2.SetProperty("name", []byte("this is vertex 2"))
	vertex2.SetProperty("nickname", []byte("node 2"))
	edge1.SetProperty("name", []byte("this is edge 1"))
	edge1.SetProperty("nickname", []byte("connector 1"))
	edge2.SetProperty("name", []byte("this is edge 2"))
	edge2.SetProperty("nickname", []byte("connector 2"))

	propval,_ = vertex1.Property("name")
	assert.Equal(t, []byte("this is vertex 1"), propval)
	propval,_ = vertex1.Property("nickname")
	assert.Equal(t, []byte("node 1"), propval)
	propval,_ = vertex2.Property("name")
	assert.Equal(t, []byte("this is vertex 2"), propval)
	propval,_ = vertex2.Property("nickname")
	assert.Equal(t, []byte("node 2"), propval)
	propval,_ = edge1.Property("name")
	assert.Equal(t, []byte("this is edge 1"), propval)
	propval,_ = edge1.Property("nickname")
	assert.Equal(t, []byte("connector 1"), propval)
	propval,_ = edge2.Property("name")
	assert.Equal(t, []byte("this is edge 2"), propval)
	propval,_ = edge2.Property("nickname")
	assert.Equal(t, []byte("connector 2"), propval)

	assert.True(t, vertex1.DelProperty("nonexistentproperty") == nil)
	assert.True(t, vertex1.DelProperty("name") == nil)
	propval,_ = vertex1.Property("name")
	assert.True(t, propval == nil)
	assert.True(t, vertex1.DelProperty("nickname") == nil)
	propval,_ = vertex1.Property("nickname")
	assert.True(t, propval == nil)

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

	if graph.Capabilities().SortedKeys() == true {
		keys,_ := vertex1.PropertyKeys()
		assert.Equal(t, propkeys, keys)
	}

	for _, key := range propkeys {
		propvertexval := key + "-vertex"
		propedgeval := key + "-edge"
		propvval,_ := vertex1.Property(key)
		propeval,_ := edge1.Property(key)
		assert.Equal(t, []byte(propvertexval), propvval)
		assert.Equal(t, []byte(propedgeval), propeval)
	}

	if graph.Capabilities().Persistent() == true {

		graph.Close()

		graph.Open()
		vertex1,_ = graph.Vertex(vid1)
		vertex2,_ = graph.Vertex(vid2)
		edge1,_ = graph.Edge(eid1)
		edge2,_ = graph.Edge(eid2)

		propval,_ = vertex1.Property("name")
		assert.True(t, propval == nil)
		propval,_ = vertex1.Property("nickname")
		assert.True(t, propval == nil)
		propval,_ = vertex2.Property("name")
		assert.Equal(t, []byte("this is vertex 2"), propval)
		propval,_ = vertex2.Property("nickname")
		assert.Equal(t, []byte("node 2"), propval)
		propval,_ = edge1.Property("name")
		assert.Equal(t, []byte("this is edge 1"), propval)
		propval,_ = edge1.Property("nickname")
		assert.Equal(t, []byte("connector 1"), propval)
		propval,_ = edge2.Property("name")
		assert.Equal(t, []byte("this is edge 2"), propval)
		propval,_ = edge2.Property("nickname")
		assert.Equal(t, []byte("connector 2"), propval)

		propkeys = append(propkeys, "name")
		propkeys = append(propkeys, "nickname")
		keys1,_ := vertex1.PropertyKeys()

		assert.Equal(t, propkeys, keys1)

	}

	graph.Close()
}
