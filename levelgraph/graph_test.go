package levelgraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
	"reflect"
	//"github.com/jmhodges/levigo"
	//"fmt"
	//"github.com/lexlapax/graveldb/core"	
)

func cleanup(dbdir string) {
	os.RemoveAll(dbdir)
}

func TestOpenGraph(t *testing.T){
	dbdir := "./testing.db"
	gdb,_ := opengraph(dbdir)
	if assert.NotNil(t, &gdb) {
		assert.Equal(t, "./testing.db", gdb.dbdir)
		assert.Equal(t, "*levigo.DB", reflect.TypeOf(gdb.meta).String())
		assert.NotNil(t, gdb.meta)
		assert.NotNil(t, gdb.elements)
		assert.NotNil(t, gdb.hs)
		assert.NotNil(t, gdb.props)
		fi, _ := os.Lstat(dbdir)
		assert.True(t, fi.IsDir(), "dbdir should be a directory")
		assert.Equal(t, "testing.db", fi.Name(), "dbdir name should match" )
		assert.Equal(t, "dbdir=./testing.db", gdb.String())
	}
	gdb.Close()
	defer cleanup(dbdir)
}

func TestAddVertex(t *testing.T) {
	id := []byte("somerandomstringid")
	dbdir := "./testing.db"
	gdb,_ := opengraph(dbdir)
	vertex, err := gdb.AddVertex(nil)
	assert.True(t, vertex == nil)
	assert.Equal(t, NilValue, err)
	vertex, err = gdb.AddVertex(id)
	if assert.True(t, vertex != nil) {
		assert.Equal(t, id, vertex.Id())
		assert.Equal(t, VertexType, vertex.Elementtype)
		assert.Equal(t, nil, err)
	}
	vertexb, errb := gdb.AddVertex(id) 
	assert.True(t, vertexb == nil)
	assert.Equal(t, KeyExists, errb )
	gdb.Close()
	defer cleanup(dbdir)
}


func TestGetVertex(t *testing.T) {
	ida := []byte("somerandomstringid")
	idb := []byte("idb")
	dbdir := "./testing.db"
	gdb,_ := opengraph(dbdir)
	vertexa, _ := gdb.AddVertex(ida)
	vertexb  := gdb.Vertex(ida)
	assert.Equal(t, vertexa, vertexb)
	assert.True(t, gdb.Vertex(idb) == nil)
	vertexc, _ := gdb.AddVertex(idb)
	vertexd := gdb.Vertex(idb)
	assert.Equal(t, vertexc, vertexd)
	gdb.Close()
	defer cleanup(dbdir)
}

func TestDelVertex(t *testing.T) {
	ida := []byte("somerandomstringid")
	dbdir := "./testing.db"
	gdb,_ := opengraph(dbdir)

	err := gdb.DelVertex(nil)
	assert.Equal(t, NilValue, err)
	err = gdb.DelVertex(new(DBVertex))
	assert.Equal(t, NilValue, err)

	vertexnull := &DBVertex{DBElement{gdb, ida, VertexType}}
	err = gdb.DelVertex(vertexnull)
	assert.Equal(t, KeyDoesNotExist, err)
	vertexa, _ := gdb.AddVertex(ida)
	assert.Equal(t, vertexa, gdb.Vertex(ida))
	err = gdb.DelVertex(vertexa)
	assert.True(t, err == nil)
	assert.True(t, gdb.Vertex(ida) == nil)
	err = gdb.DelVertex(vertexa)
	assert.Equal(t, KeyDoesNotExist, err)
	gdb.Close()
	defer cleanup(dbdir)
}

func TestVertexCount(t *testing.T) {
	dbdir := "./testing.db"
	gdb,_ := opengraph(dbdir)

	assert.Equal(t, uint64(0), gdb.VertexCount())

	ida := []byte("somerandomstringid")
	vertexa,_ := gdb.AddVertex(ida)
	assert.Equal(t, uint64(1), gdb.VertexCount())

	testvertii := []*DBVertex{}
	var vertex *DBVertex
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _ = gdb.AddVertex([]byte(a + "-" + n))
			testvertii = append(testvertii, vertex)
		}
	}
	numv := len(testvertii)
	assert.Equal(t, uint64(numv + 1), gdb.VertexCount())
	gdb.DelVertex(vertexa)
	assert.Equal(t, uint64(numv), gdb.VertexCount())
	for i :=0; i < numv; i++ {
		gdb.DelVertex(testvertii[i])
		assert.Equal(t, uint64(numv - (i + 1)), gdb.VertexCount() )
	}
	assert.Equal(t, uint64(0), gdb.VertexCount())
	gdb.Close()
	defer cleanup(dbdir)
}


func TestVertices(t *testing.T) {
	ida := []byte("somerandomstringid")
	dbdir := "./testing.db"
	gdb,_ := opengraph(dbdir)
	testvertii := []*DBVertex{}
	var vertex *DBVertex
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	for _,a := range alpha {
		for _,n := range numb { 
			vertex, _= gdb.AddVertex([]byte(a + "-" + n))
			testvertii = append(testvertii, vertex)
		}
	}
	assert.Equal(t, testvertii, gdb.Vertices())
	lastvertex, _ := gdb.AddVertex(ida)
	assert.NotEqual(t, testvertii, gdb.Vertices())
	//keys are lexicaly ordered.. lastvertex should be the last in the list
	testvertii = gdb.Vertices()
	assert.Equal(t, lastvertex, testvertii[len(testvertii) - 1])
	gdb.Close()
	defer cleanup(dbdir)
}