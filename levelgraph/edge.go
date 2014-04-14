package levelgraph


import (
		//"bytes"
		//"fmt"
		//"errors"
		//"os"
		//"github.com/jmhodges/levigo"
		"github.com/lexlapax/graveldb/core"
)

type DBEdge struct {
	DBElement
	label string
}


func (edge *DBEdge) Label() (string) {
	return edge.label
}

func (edge *DBEdge) VertexOut() (core.Vertex) {
	return new(DBVertex)
}

func (edge *DBEdge) VertexIn() (core.Vertex) {
	return new(DBVertex)
}

func (edge *DBEdge) String() (string) {
	return ""
}
