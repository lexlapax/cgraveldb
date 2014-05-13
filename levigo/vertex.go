package levigo


import (
		"fmt"
		"github.com/lexlapax/graveldb/core"
)

type VertexLevigo struct {
	*AtomLevigo
}


func (vertex *VertexLevigo) IterEdges(direction core.Direction, labels ...string) <-chan core.Edge {
	return vertex.db.vertexIterEdges(direction, vertex, labels...)
}

func (vertex *VertexLevigo) IterOutEdges(labels ...string) <-chan core.Edge {
	return vertex.IterEdges(core.DirOut, labels...)
}

func (vertex *VertexLevigo) IterInEdges(labels ...string) <-chan core.Edge {
	return vertex.IterEdges(core.DirIn, labels...)
}

func (vertex *VertexLevigo) IterVertices(direction core.Direction, labels ...string) <-chan core.Vertex {
	return vertex.db.vertexIterVertices(direction, vertex, labels...)
}

func (vertex *VertexLevigo) Vertices(direction core.Direction, labels ...string) ([]core.Vertex, error) {
	return vertex.db.vertexVertices(direction, vertex, labels...)
}

func (vertex *VertexLevigo) Edges(direction core.Direction, labels ...string) ([]core.Edge, error) {
	return vertex.db.vertexEdges(direction, vertex, labels...)
}

func (vertex *VertexLevigo) OutEdges(labels ...string) ([]core.Edge, error) {
	return vertex.db.vertexEdges(core.DirOut, vertex, labels...) 
}

func (vertex *VertexLevigo) InEdges(labels ...string) ([]core.Edge, error) {
	return vertex.db.vertexEdges(core.DirIn, vertex, labels...) 
}

func (vertex *VertexLevigo) AddEdge(id string, invertex core.Vertex, label string) (core.Edge, error) {
	return vertex.db.AddEdge(id, vertex, invertex, label)
}

func (vertex *VertexLevigo) Query() core.QueryVertex {
	return core.NewBaseVertexQuery(vertex)
}

func (vertex *VertexLevigo) String() (string) {
	str := fmt.Sprintf("<VertexLevigo:%v@%v>",vertex.Id(), vertex.db)
	return str
}

