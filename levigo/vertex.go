package levigo


import (
		"fmt"
		//"github.com/lexlapax/graveldb/core"
)

type VertexLevigo struct {
	*AtomLevigo
}

/*

func (vertex *VertexLevigo) Vertices() ([]*VertexLevigo) {
	return nil
}

*/

func (vertex *VertexLevigo) OutEdges() ([]*EdgeLevigo) {
	return vertex.db.VertexOutEdges(vertex)
}

func (vertex *VertexLevigo) InEdges() ([]*EdgeLevigo) {
	return vertex.db.VertexInEdges(vertex)
}

func (vertex *VertexLevigo) AddEdge(id []byte, invertex *VertexLevigo, label string) (*EdgeLevigo, error) {
	return vertex.db.AddEdge(id, vertex, invertex, label)
}


func (vertex *VertexLevigo) String() (string) {
	str := fmt.Sprintf("<VertexLevigo:%v@%v>",vertex.IdAsString(), vertex.db)
	return str
}
