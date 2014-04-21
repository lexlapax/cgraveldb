package levelgraph


import (
		"fmt"
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

func (vertex *DBVertex) AddEdge(id []byte, invertex *DBVertex, label string) (*DBEdge, error) {
	return vertex.db.AddEdge(id, vertex, invertex, label)
}


func (vertex *DBVertex) String() (string) {
	str := fmt.Sprintf("<DBVertex:%v@%v>",vertex.IdAsString(), vertex.db)
	return str
}
