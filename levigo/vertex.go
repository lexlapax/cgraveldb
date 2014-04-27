package levigo


import (
		"fmt"
		"github.com/lexlapax/graveldb/core"
)

type VertexLevigo struct {
	*AtomLevigo
}


func (vertex *VertexLevigo) Vertices(direction core.Direction, labels ...string) ([]core.Vertex, error) {
	return nil, nil
}

func (vertex *VertexLevigo) Edges(direction core.Direction, labels ...string) ([]core.Edge, error) {
	return nil, nil 
}

func (vertex *VertexLevigo) OutEdges(labels ...string) ([]core.Edge, error) {
	return vertex.db.VertexOutEdges(vertex)
}

func (vertex *VertexLevigo) InEdges(labels ...string) ([]core.Edge, error) {
	return vertex.db.VertexInEdges(vertex)
}

func (vertex *VertexLevigo) AddEdge(id []byte, invertex core.Vertex, label string) (core.Edge, error) {
	return vertex.db.AddEdge(id, vertex, invertex, label)
}


func (vertex *VertexLevigo) String() (string) {
	str := fmt.Sprintf("<VertexLevigo:%v@%v>",vertex.IdAsString(), vertex.db)
	return str
}

