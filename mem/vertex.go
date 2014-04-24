package mem

import (
		"github.com/lexlapax/graveldb/core"
		//mapset "github.com/deckarep/golang-set"
		//"fmt"
)

type VertexMem struct {
	*AtomMem
	// outedges map[string]mapset.Set
	// inedges map[string]mapset.Set
	outedges map[string]*edgeSet
	inedges map[string]*edgeSet
}

func NewVertexMem(db *GraphMem, id []byte) *VertexMem {
	//vertex := &VertexMem{NewAtomMem(db, id, VertexType), make(map[string]mapset.Set), make(map[string]mapset.Set)}
	vertex := &VertexMem{NewAtomMem(db, id, VertexType), make(map[string]*edgeSet), make(map[string]*edgeSet)}
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
		//fmt.Printf("forward edges=%v\n",forward)
		if err != nil {return []core.Edge{}, err}
		reverse, err := vertex.InEdges(labels...)
		//fmt.Printf("reverse edges=%v\n",reverse)
		if err != nil {return []core.Edge{}, err}
		return append(forward, reverse...), nil
	}
	
}

func (vertex *VertexMem) Vertices(direction core.Direction, labels ...string) ([]core.Vertex, error) {
	return nil, nil
}

func iterEdgeSet(edgechan <-chan interface{}, edges *[]core.Edge) {
	for s:= range edgechan {
		*edges = append(*edges, s.(core.Edge))
	}
}

func (vertex *VertexMem) OutEdges(labels ...string) ([]core.Edge, error) {
	totaledges := []core.Edge{}
	if len(labels) == 0 {
		for _, edgeset := range vertex.outedges {
			if edgeset != nil && edgeset.count() > 0 {
				totaledges = append(totaledges, edgeset.members()...)
			}
		}
	} else {
		for _, label := range labels {
			if edgeset, ok := vertex.outedges[label]; ok {
				if edgeset != nil && edgeset.count() > 0 {
					totaledges = append(totaledges, edgeset.members()...)
				}
			}
		}
	}

	//fmt.Printf("totaledges=%v\n", totaledges)
	return totaledges, nil
}

func (vertex *VertexMem) InEdges(labels ...string) ([]core.Edge, error) {
	totaledges := []core.Edge{}
	if len(labels) == 0 {
		for _, edgeset := range vertex.inedges {
			if edgeset != nil && edgeset.count() > 0 {
				totaledges = append(totaledges, edgeset.members()...)
			}
		}
	} else {
		for _, label := range labels {
			if edgeset, ok := vertex.inedges[label]; ok {
				if edgeset != nil && edgeset.count() > 0 {
					totaledges = append(totaledges, edgeset.members()...)
				}
			}
		}
	}

	//fmt.Printf("totaledges=%v\n", totaledges)
	return totaledges, nil
}

func (vertex *VertexMem) AddEdge(id []byte, invertex core.Vertex, label string) (core.Edge, error) {
	return vertex.graph.AddEdge(id, vertex, invertex, label)
}

func (vertex *VertexMem) addOutEdge(edge core.Edge) {
	edges := vertex.outedges[edge.Label()]
	if edges == nil {
		edges = newEdgeSet()
	}
	edges.add(edge)
	vertex.outedges[edge.Label()] = edges
}

func (vertex *VertexMem) addInEdge(edge core.Edge) {
	edges := vertex.inedges[edge.Label()]
	if edges == nil {
		edges = newEdgeSet()
	}
	edges.add(edge)
	vertex.inedges[edge.Label()] = edges
}

func (vertex *VertexMem) delOutEdge(edge core.Edge) {
	edges := vertex.outedges[edge.Label()]
	if edges == nil {
		edges = newEdgeSet()
	}
	edges.del(edge)
	vertex.outedges[edge.Label()] = edges
}

func (vertex *VertexMem) delInEdge(edge core.Edge) {
	edges := vertex.inedges[edge.Label()]
	if edges == nil {
		edges = newEdgeSet()
	}
	edges.del(edge)
	vertex.inedges[edge.Label()] = edges
}
