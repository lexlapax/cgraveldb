package mem

import (
	"github.com/lexlapax/graveldb/core"
	"github.com/stretchr/testify/assert"
	"testing"
	//"os"
	//"reflect"
	//"github.com/jmhodges/levigo"
	//"fmt"
	//"github.com/lexlapax/graveldb/core"	
)

var graphimp = "mem"

func TestGraphEmpty(t *testing.T){
	gdb := core.GetGraph(graphimp)
	gdb.Open()
	defer gdb.Close()
	vertices, _ := gdb.Vertices()
	assert.Equal(t, 0, len(vertices))
	edges, _ := gdb.Edges()
	assert.Equal(t, 0, len(edges))
}


func TestGraphAdd(t *testing.T) {

}