package levelgraph


import (
		//"bytes"
		//"fmt"
		//"errors"
		//"os"
		//"github.com/jmhodges/levigo"
		"github.com/lexlapax/graveldb/core"
)

type DBVertex struct {
	DBElement
}

func (vertex *DBVertex) Vertices() ([]core.Vertex) {
	return nil
}

func (vertex *DBVertex) OutEdges() ([]core.Edge) {
	return nil
}

func (vertex *DBVertex) InEdges() ([]core.Edge) {
	return nil
}

func (vertex *DBVertex) AddEdge(outvertex core.Vertex, label string) (core.Edge, error) {
	return new(DBEdge), nil
}

func (vertex *DBVertex) String() (string) {
	return ""
}
