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

var graphimp = "mem"
func init() {
	mem.Register()
}

func TestGraphEmpty(t *testing.T){
	graph := core.GetGraph(graphimp)
	graph.Open()
	defer graph.Close()

	//todo - graph capabilites based 
	vertices, _ := graph.Vertices()
	assert.Equal(t, 0, len(vertices))
	edges, _ := graph.Edges()
	assert.Equal(t, 0, len(edges))
}


func TestGraphAdd(t *testing.T) {
	graph := core.GetGraph(graphimp)
	graph.Open()
	defer graph.Close()

	a,_ := graph.AddVertex(nil)
	b,_ := graph.AddVertex(nil)
    edge, _ := graph.AddEdge(nil, a, b, "knows")

	//todo - graph capabilites based 
	vertices, _ := graph.Vertices()
	assert.Equal(t, 2, len(vertices))
	edges, _ := graph.Edges()
	assert.Equal(t, 1, len(edges))

	graph.DelVertex(a)

	vertices, _ = graph.Vertices()
	assert.Equal(t, 1, len(vertices))
	edges, _ = graph.Edges()
	assert.Equal(t, 0, len(edges))
	err := graph.DelEdge(edge)
	assert.True(t, err != nil)
}

func TestGraphSetProperties(t *testing.T) {
	graph := core.GetGraph(graphimp)
	graph.Open()
	defer graph.Close()
	//todo graph capabilities based
	a,_ := graph.AddVertex(nil)
	b,_ := graph.AddVertex(nil)
	graph.AddEdge(nil, a, b, "knows")
	graph.AddEdge(nil, a, b, "knows")
	edges, _ := b.Edges(core.DirIn)
    for _,edge := range edges {
            edge.SetProperty("key", []byte("value"))
    }
}

func TestNodesWithKeyValue(t *testing.T) {
	t.Skip()
	//todo
}

