package levigo


import (
		"fmt"
		"github.com/lexlapax/graveldb/core"
)

type VertexLevigo struct {
	*AtomLevigo
}


func (vertex *VertexLevigo) Vertices(direction core.Direction, labels ...string) ([]core.Vertex, error) {
	var forward, reverse []core.Vertex
	var err error
	if direction == core.DirOut {
		forward, err = vertex.db.vertexVertices(core.DirOut, vertex, labels...)
		return forward, err
	} else if direction == core.DirIn {
		reverse, err = vertex.db.vertexVertices(core.DirIn, vertex, labels...)
		return reverse, err
	} else {
		forward, err = vertex.db.vertexVertices(core.DirOut, vertex, labels...)
		//fmt.Printf("forward vertices=%v,err=%v\n",forward, err)
		if err != nil {return []core.Vertex{}, err}
		reverse, err = vertex.db.vertexVertices(core.DirIn, vertex, labels...)
		//fmt.Printf("reverse vertices=%v\n",reverse)
		if err != nil {return []core.Vertex{}, err}
		return append(forward, reverse...), nil
	}

	return nil, nil 
}

func (vertex *VertexLevigo) Edges(direction core.Direction, labels ...string) ([]core.Edge, error) {
	var forward, reverse []core.Edge
	var err error
	if direction == core.DirOut {
		forward, err = vertex.OutEdges(labels...)
		return forward, err
	} else if direction == core.DirIn {
		reverse, err = vertex.InEdges(labels...)
		return reverse, err
	} else {
		forward, err := vertex.OutEdges(labels...)
		//fmt.Printf("forward edges=%v\n",forward)
		if err != nil {return []core.Edge{}, err}
		reverse, err := vertex.InEdges(labels...)
		//fmt.Printf("reverse edges=%v\n",reverse)
		if err != nil {return []core.Edge{}, err}
		return append(forward, reverse...), nil
	}

	return nil, nil 
}

func (vertex *VertexLevigo) OutEdges(labels ...string) ([]core.Edge, error) {
	return vertex.db.VertexOutEdges(vertex, labels...)
}

func (vertex *VertexLevigo) InEdges(labels ...string) ([]core.Edge, error) {
	return vertex.db.VertexInEdges(vertex, labels...)
}

func (vertex *VertexLevigo) AddEdge(id string, invertex core.Vertex, label string) (core.Edge, error) {
	return vertex.db.AddEdge(id, vertex, invertex, label)
}


func (vertex *VertexLevigo) String() (string) {
	str := fmt.Sprintf("<VertexLevigo:%v@%v>",vertex.Id(), vertex.db)
	return str
}

