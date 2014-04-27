package coretest

import (
	"github.com/stretchr/testify/assert"
)

type AtomTestSuite struct {
	BaseTestSuite
}

func (suite *AtomTestSuite) TestAtomProperty(){
	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	
	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)

	eid1 := []byte("edge1")
	eid2 := []byte("edge2")

	edge1, _ := suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := suite.TestGraph.AddEdge(eid2, vertex2, vertex1, "2 to 1")

	propval,_ := vertex1.Property("name")
	assert.True(suite.T(), propval == nil)
	propval,_ = vertex2.Property("name")
	assert.True(suite.T(),  propval == nil)
	propval,_ = edge1.Property("name")
	assert.True(suite.T(), propval == nil)
	propval,_ = edge2.Property("name")
	assert.True(suite.T(), propval == nil)

	vertex1.SetProperty("name", nil)
	propval,_ = vertex1.Property("name")
	assert.Equal(suite.T(), []byte{}, propval)	
	vertex1.SetProperty("name", []byte("some value"))
	propval,_ = vertex1.Property("name")
	assert.Equal(suite.T(), []byte("some value"), propval)


	vertex1.SetProperty("name", []byte("this is vertex 1"))
	vertex1.SetProperty("nickname", []byte("node 1"))
	vertex2.SetProperty("name", []byte("this is vertex 2"))
	vertex2.SetProperty("nickname", []byte("node 2"))
	edge1.SetProperty("name", []byte("this is edge 1"))
	edge1.SetProperty("nickname", []byte("connector 1"))
	edge2.SetProperty("name", []byte("this is edge 2"))
	edge2.SetProperty("nickname", []byte("connector 2"))

	propval,_ = vertex1.Property("name")
	assert.Equal(suite.T(), []byte("this is vertex 1"), propval)
	propval,_ = vertex1.Property("nickname")
	assert.Equal(suite.T(), []byte("node 1"), propval)
	propval,_ = vertex2.Property("name")
	assert.Equal(suite.T(), []byte("this is vertex 2"), propval)
	propval,_ = vertex2.Property("nickname")
	assert.Equal(suite.T(), []byte("node 2"), propval)
	propval,_ = edge1.Property("name")
	assert.Equal(suite.T(), []byte("this is edge 1"), propval)
	propval,_ = edge1.Property("nickname")
	assert.Equal(suite.T(), []byte("connector 1"), propval)
	propval,_ = edge2.Property("name")
	assert.Equal(suite.T(), []byte("this is edge 2"), propval)
	propval,_ = edge2.Property("nickname")
	assert.Equal(suite.T(), []byte("connector 2"), propval)

	assert.True(suite.T(), vertex1.DelProperty("nonexistentproperty") == nil)
	assert.True(suite.T(), vertex1.DelProperty("name") == nil)
	propval,_ = vertex1.Property("name")
	assert.True(suite.T(), propval == nil)
	assert.True(suite.T(), vertex1.DelProperty("nickname") == nil)
	propval,_ = vertex1.Property("nickname")
	assert.True(suite.T(), propval == nil)

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

	if suite.TestGraph.Capabilities().SortedKeys() == true {
		keys,_ := vertex1.PropertyKeys()
		assert.Equal(suite.T(), propkeys, keys)
	}

	for _, key := range propkeys {
		propvertexval := key + "-vertex"
		propedgeval := key + "-edge"
		propvval,_ := vertex1.Property(key)
		propeval,_ := edge1.Property(key)
		assert.Equal(suite.T(), []byte(propvertexval), propvval)
		assert.Equal(suite.T(), []byte(propedgeval), propeval)
	}

	if suite.TestGraph.Capabilities().Persistent() == true {

		suite.TestGraph.Close()

		suite.TestGraph.Open()
		vertex1,_ = suite.TestGraph.Vertex(vid1)
		vertex2,_ = suite.TestGraph.Vertex(vid2)
		edge1,_ = suite.TestGraph.Edge(eid1)
		edge2,_ = suite.TestGraph.Edge(eid2)

		propval,_ = vertex1.Property("name")
		assert.True(suite.T(), propval == nil)
		propval,_ = vertex1.Property("nickname")
		assert.True(suite.T(), propval == nil)
		propval,_ = vertex2.Property("name")
		assert.Equal(suite.T(), []byte("this is vertex 2"), propval)
		propval,_ = vertex2.Property("nickname")
		assert.Equal(suite.T(), []byte("node 2"), propval)
		propval,_ = edge1.Property("name")
		assert.Equal(suite.T(), []byte("this is edge 1"), propval)
		propval,_ = edge1.Property("nickname")
		assert.Equal(suite.T(), []byte("connector 1"), propval)
		propval,_ = edge2.Property("name")
		assert.Equal(suite.T(), []byte("this is edge 2"), propval)
		propval,_ = edge2.Property("nickname")
		assert.Equal(suite.T(), []byte("connector 2"), propval)

		//propkeys = append(propkeys, "name")
		//propkeys = append(propkeys, "nickname")
		keys1,_ := vertex1.PropertyKeys()

		assert.Equal(suite.T(), propkeys, keys1)

	}
}
