package mem

import (
		"github.com/lexlapax/graveldb/core"
		mapset "github.com/deckarep/golang-set"
)

type VertexMem struct {
	*AtomMem
	outedges map[string]mapset.Set
	inedges map[string]mapset.Set
}

func NewVertexMem(db *GraphMem, id []byte) *VertexMem {
	vertex := &VertexMem{NewAtomMem(db, id, VertexType), make(map[string]mapset.Set), make(map[string]mapset.Set)}
	return vertex
}


func (vertex *VertexMem) Edges(direction core.Direction, labels ...string) ([]core.Edge, error) {
	var forward, reverse []core.Edge
	var err error
	if direction == core.DirForward {
		forward, err = vertex.OutEdges(labels...)
		return forward, err
	} else if direction == core.DirReverse {
		reverse, err = vertex.InEdges(labels...)
		return reverse, err
	} else {
		forward, err := vertex.OutEdges(labels...)
		if err != nil {return []core.Edge{}, err}
		reverse, err := vertex.InEdges(labels...)
		if err != nil {return []core.Edge{}, err}
		return append(forward, reverse...), nil
	}
	
}

func (vertex *VertexMem) Vertices(direction core.Direction, labels ...string) ([]core.Vertex, error) {
	return nil, nil
}

func (vertex *VertexMem) OutEdges(labels ...string) ([]core.Edge, error) {
	totaledges := []core.Edge{}
	for _, label := range labels {
		edges := vertex.outedges[label]
		if edges != nil {
			cs := make(chan core.Edge)
			for s := range cs {
        		totaledges = append(totaledges, s)
    		}
		}
	}
	return totaledges, nil
}

func (vertex *VertexMem) InEdges(labels ...string) ([]core.Edge, error) {
	totaledges := []core.Edge{}
	for _, label := range labels {
		edges := vertex.inedges[label]
		if edges != nil {
			cs := make(chan core.Edge)
			for s := range cs {
        		totaledges = append(totaledges, s)
    		}
		}
	}
	return totaledges, nil
}

func (vertex *VertexMem) AddEdge(id []byte, invertex core.Vertex, label string) (core.Edge, error) {
	return vertex.graph.AddEdge(id, vertex, invertex, label)
}

func (vertex *VertexMem) addOutEdge(edge core.Edge) {
	edges := vertex.outedges[edge.Label()]
	if edges == nil {
		edges = mapset.NewSet()
	}
	edges.Add(edge)
	vertex.outedges[edge.Label()] = edges
}

func (vertex *VertexMem) addInEdge(edge core.Edge) {
	edges := vertex.inedges[edge.Label()]
	if edges == nil {
		edges = mapset.NewSet()
	}
	edges.Add(edge)
	vertex.inedges[edge.Label()] = edges
}

func (vertex *VertexMem) delOutEdge(edge core.Edge) {
	edges := vertex.outedges[edge.Label()]
	if edges == nil {
		edges = mapset.NewSet()
	}
	edges.Remove(edge)
	vertex.outedges[edge.Label()] = edges
}

func (vertex *VertexMem) delInEdge(edge core.Edge) {
	edges := vertex.inedges[edge.Label()]
	if edges == nil {
		edges = mapset.NewSet()
	}
	edges.Remove(edge)
	vertex.inedges[edge.Label()] = edges
}
