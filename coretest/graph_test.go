package coretest

import (
	"github.com/lexlapax/graveldb/core"
	"github.com/stretchr/testify/assert"
	"testing"
	//"os"
	//"reflect"
	//"github.com/jmhodges/levigo"
	"fmt"
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
	fmt.Printf("before add vertex\n")

	a,_ := graph.AddVertex(nil)
	b,_ := graph.AddVertex(nil)
    edge, _ := graph.AddEdge(nil, a, b, "knows")

	//todo - graph capabilites based 
	vertices, _ := graph.Vertices()
	assert.Equal(t, 2, len(vertices))
	edges, _ := graph.Edges()
	assert.Equal(t, 1, len(edges))

	fmt.Printf("before delvertex\n")
	graph.DelVertex(a)
	vertices, _ = graph.Vertices()
	assert.Equal(t, 1, len(vertices))
	edges, _ = graph.Edges()
	assert.Equal(t, 0, len(edges))
	fmt.Printf("before deledge\n")
	err := graph.DelEdge(edge)
	assert.True(t, err != nil)
}
