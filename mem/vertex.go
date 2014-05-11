package mem

import (
		"github.com/lexlapax/graveldb/core"
		//mapset "github.com/deckarep/golang-set"
		//"fmt"
)

type VertexMem struct {
	*AtomMem
	outedges map[string]*core.AtomSet
	inedges map[string]*core.AtomSet
}

func NewVertexMem(db *GraphMem, id string) *VertexMem {
	//vertex := &VertexMem{NewAtomMem(db, id, VertexType), make(map[string]mapset.Set), make(map[string]mapset.Set)}
	vertex := &VertexMem{NewAtomMem(db, id, core.VertexType), make(map[string]*core.AtomSet), make(map[string]*core.AtomSet)}
	return vertex
}

func getEdges(ch chan core.Edge, edgemap map[string]*core.AtomSet, labels ...string) {
	if len(labels) == 0 { 
		for _, edgeset := range edgemap {
			if edgeset != nil && edgeset.Count() > 0 {
				for _, edge := range edgeset.Members() {
					ch <- edge.(core.Edge)
				}
			}
		}
	} else {
		for _, label := range labels {
			if edgeset, ok := edgemap[label]; ok {
				if edgeset != nil && edgeset.Count() > 0 {
					for _, edge := range edgeset.Members() {
						ch <- edge.(core.Edge)
					}
				}
			}
		}

	}
}

func (vertex *VertexMem) Edges(direction core.Direction, labels ...string) ([]core.Edge, error) {
	var edges []core.Edge
	for edge := range vertex.IterEdges(direction, labels...) {
		edges = append(edges, edge)
	}
	return edges, nil
}


func (vertex *VertexMem) IterEdges(direction core.Direction, labels ...string) <-chan core.Edge {
	ch := make(chan core.Edge)

	go func() {
		if direction == core.DirOut {
			getEdges(ch, vertex.outedges, labels...)
		} else if direction == core.DirIn {
			getEdges(ch, vertex.inedges, labels...)
		} else {
			getEdges(ch, vertex.outedges, labels...)
			getEdges(ch, vertex.inedges, labels...)
		}
		close(ch)
	}()
	return ch
}


func (vertex *VertexMem) IterOutEdges(labels ...string) <-chan core.Edge {
	return vertex.IterEdges(core.DirOut, labels...)
}

func (vertex *VertexMem) IterInEdges(labels ...string) <-chan core.Edge {
	return vertex.IterEdges(core.DirIn, labels...)
}

func (vertex *VertexMem) OutEdges(labels ...string) ([]core.Edge, error) {
	totaledges := []core.Edge{}

	for edge := range vertex.IterEdges(core.DirOut, labels...) {
		totaledges = append(totaledges, edge)
	}

	//fmt.Printf("totaledges=%v\n", totaledges)
	return totaledges, nil
}


func (vertex *VertexMem) InEdges(labels ...string) ([]core.Edge, error) {
	totaledges := []core.Edge{}
	for edge := range vertex.IterEdges(core.DirIn, labels...) {
		totaledges = append(totaledges, edge)
	}
	//fmt.Printf("totaledges=%v\n", totaledges)
	return totaledges, nil
}

func (vertex *VertexMem) IterVertices(direction core.Direction, labels ...string) <-chan core.Vertex {
	ch := make(chan core.Vertex)

	go func() {
		if direction == core.DirOut {
			for edge := range vertex.IterOutEdges(labels...) {
				vertex,_ := edge.VertexIn()
				ch <- vertex
			}
		} else if direction == core.DirIn {
			for edge := range vertex.IterInEdges(labels...) {
				vertex,_ := edge.VertexOut()
				ch <- vertex
			}
		} else {
			for edge := range vertex.IterOutEdges(labels...) {
				vertex,_ := edge.VertexIn()
				ch <- vertex
			}
			for edge := range vertex.IterInEdges(labels...) {
				vertex,_ := edge.VertexOut()
				ch <- vertex
			}
		}
		close(ch)
		}()
	return ch
}

func (vertex *VertexMem) Vertices(direction core.Direction, labels ...string) ([]core.Vertex, error) {
	var vertices []core.Vertex
	for tmpvertex := range vertex.IterVertices(direction, labels...) {
		vertices = append(vertices, tmpvertex)
	}
	return vertices, nil

}

func iterEdgeSet(edgechan <-chan interface{}, edges *[]core.Edge) {
	for s:= range edgechan {
		*edges = append(*edges, s.(core.Edge))
	}
}

func (vertex *VertexMem) OutVertices(labels ...string) ([]core.Vertex, error) {
	return vertex.Vertices(core.DirOut, labels...)
}

func (vertex *VertexMem) InVertices(labels ...string) ([]core.Vertex, error) {
	return vertex.Vertices(core.DirIn, labels...)
}


func (vertex *VertexMem) AddEdge(id string, invertex core.Vertex, label string) (core.Edge, error) {
	return vertex.graph.AddEdge(id, vertex, invertex, label)
}

func (vertex *VertexMem) addOutEdge(edge core.Edge) {
	edges := vertex.outedges[edge.Label()]
	if edges == nil {
		edges = core.NewAtomSet()
	}
	edges.Add(edge)
	vertex.outedges[edge.Label()] = edges
}

func (vertex *VertexMem) addInEdge(edge core.Edge) {
	edges := vertex.inedges[edge.Label()]
	if edges == nil {
		edges = core.NewAtomSet()
	}
	edges.Add(edge)
	vertex.inedges[edge.Label()] = edges
}

func (vertex *VertexMem) delOutEdge(edge core.Edge) {
	edges := vertex.outedges[edge.Label()]
	if edges == nil {
		edges = core.NewAtomSet()
	}
	edges.Del(edge)
	vertex.outedges[edge.Label()] = edges
}

func (vertex *VertexMem) delInEdge(edge core.Edge) {
	edges := vertex.inedges[edge.Label()]
	if edges == nil {
		edges = core.NewAtomSet()
	}
	edges.Del(edge)
	vertex.inedges[edge.Label()] = edges
}
