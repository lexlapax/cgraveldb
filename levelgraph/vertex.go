package levelgraph


import (
		//"bytes"
		"fmt"
		//"errors"
		//"os"
		//"github.com/jmhodges/levigo"
		//"github.com/lexlapax/graveldb/core"
)

type DBVertex struct {
	*DBElement
}

func (vertex *DBVertex) Vertices() ([]*DBVertex) {
	return nil
}

func (vertex *DBVertex) OutEdges() ([]*DBEdge) {
	return nil
}

func (vertex *DBVertex) InEdges() ([]*DBEdge) {
	return nil
}

func (vertex *DBVertex) AddEdge(outvertex *DBVertex, label string) (*DBEdge, error) {
	return nil, nil
}


func (vertex *DBVertex) String() (string) {
	str := fmt.Sprintf("<DBVertex:%v@%v>",vertex.IdAsString(), vertex.db)
	return str
}
