package mem

import (
		"github.com/lexlapax/graveldb/core"
		"sync"
		"errors"
)

const (
	DriverName                = "mem"
	)

var (
	ErrAlreadyExists = errors.New("the object with the id already exists")
)

func init() {
	core.Register(DriverName, &GraphMem{})
	// core.RegisterGraphDb(DriverName, func(args ...interface{}) (*GraphMem, error) {
	// 	graph, err := NewGraph()
	// 	return graph, nil
	// })
}

type GraphMem struct {
	vertices map[string]*VertexMem 
	edges map[string]*EdgeMem
	vertexlock *sync.RWMutex 
	edgelock *sync.RWMutex 
}

func (graph *GraphMem) AddVertex(id []byte) (core.Vertex, error) {
	graph.vertexlock.Lock()
	defer graph.vertexlock.Unlock()
	vertex := NewVertexMem(graph, id)
	graph.vertices[string(id[:])] = vertex
	return vertex, nil
}

func (graph *GraphMem) Vertex(id []byte) (core.Vertex, error) {
	graph.vertexlock.RLock()
	defer graph.vertexlock.RUnlock()
 	if val, ok := graph.vertices[string(id[:])]; ok {
        return val, nil
    }
	return nil, nil
}

func (graph *GraphMem) DelVertex(vertex core.Vertex) error {
	graph.vertexlock.Lock()
	defer graph.vertexlock.Unlock()
	edges, _ := vertex.Edges(core.DirAny)
	for _, edge := range edges {
		graph.DelEdge(edge)
	}
	delete(graph.vertices, string(vertex.Id()[:]))
	return nil
}

func (graph *GraphMem) Vertices() ([]core.Vertex, error) {
	graph.vertexlock.RLock()
	defer graph.vertexlock.RUnlock()
	vertices := []core.Vertex{}
	for _, v := range graph.vertices {
		vertices = append(vertices, v)
	}
	return vertices, nil
}

func (graph *GraphMem) AddEdge(id []byte, outvertex core.Vertex, invertex core.Vertex, label string) (core.Edge, error) {
	if _, ok := graph.edges[string(id[:])]; ok {
		return nil, ErrAlreadyExists
	}
	subject := outvertex.(*VertexMem)
	object := invertex.(*VertexMem)
	edge :=  NewEdgeMem(graph, id, subject, object, label)
	graph.edges[string(id[:])] = edge
	subject.addOutEdge(edge)
	object.addInEdge(edge)
	return edge, nil
}

func (graph *GraphMem) Edge(id []byte) (core.Edge, error) {
	graph.edgelock.RLock()
	defer graph.edgelock.RUnlock()
 	if val, ok := graph.edges[string(id[:])]; ok {
        return val, nil
    }
	return nil, nil
}

func (graph *GraphMem) DelEdge(edge core.Edge) error {
	v, _ := edge.VertexOut()
	vertexout := v.(*VertexMem)
	v, _ = edge.VertexOut()
	vertexin := v.(*VertexMem)
	vertexout.delOutEdge(edge)
	vertexin.delOutEdge(edge)
	return nil
}

func (graph *GraphMem) Edges() ([]core.Edge, error) {
	return nil, nil
}

func (graph *GraphMem) EdgeCount() uint {
	return uint(len(graph.edges))
}

func (graph *GraphMem) VertexCount() uint {
	return uint(len(graph.vertices))
}

func (graph *GraphMem) Open() error {
	graph.vertices = make(map[string]*VertexMem)
	graph.edges = make(map[string]*EdgeMem)
	graph.vertexlock = &sync.RWMutex{}
	graph.edgelock = &sync.RWMutex{}
	return nil
}

func (graph *GraphMem) Close() error {
	graph.vertices = nil
	graph.edges = nil
	graph.vertexlock = nil
	graph.edgelock = nil
	return nil
}
