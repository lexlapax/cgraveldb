package mem

import (
		"github.com/lexlapax/graveldb/core"
		"sync"
		"errors"
		"fmt"
		"strconv"
)

const (
	GraphImpl                = "mem"
	)

var (
	ErrDoesntExist = errors.New("the object with the id does not exist")
	ErrAlreadyExists = errors.New("the object with the id already exists")
	ErrNilValue = errors.New("value passed cannot be nil")
	register sync.Once
)

func Register() {
	register.Do(func() {core.Register(GraphImpl, &GraphMem{})} )
}

func init() {
	Register()
}

func NewGraphMem() core.Graph {
	graph := &GraphMem{}
	graph.Open()
	return graph
}

type GraphMem struct {
	vertices map[string]*VertexMem 
	edges map[string]*EdgeMem
	vertexlock *sync.RWMutex 
	edgelock *sync.RWMutex
	nextid uint64
	isopen bool 
}

func (graph *GraphMem) AddVertex(id []byte) (core.Vertex, error) {
	var idstr []byte
	if id == nil {
		idstr = graph.generateId()
	} else {
		var vok, eok bool
		_, eok = graph.edges[string(id[:])]
		_, vok = graph.vertices[string(id[:])]
		if vok || eok {
			return nil, ErrAlreadyExists
		}
		idstr = id
	}
	graph.vertexlock.Lock()
	defer graph.vertexlock.Unlock()
	vertex := NewVertexMem(graph, idstr)
	graph.vertices[string(idstr[:])] = vertex
	return vertex, nil
}

func (graph *GraphMem) Vertex(id []byte) (core.Vertex, error) {
	if id == nil { return nil, ErrNilValue}
	graph.vertexlock.RLock()
	defer graph.vertexlock.RUnlock()
 	if val, ok := graph.vertices[string(id[:])]; ok {
        return val, nil
    }
	return nil, nil
}

func (graph *GraphMem) DelVertex(vertex core.Vertex) error {
	if vertex == nil { return ErrNilValue}
	graph.vertexlock.Lock()
	defer graph.vertexlock.Unlock()
	edges, _ := vertex.Edges(core.DirAny)
	fmt.Printf("edges=%v\n",edges)
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
	var idstr []byte
	if id == nil {
		idstr = graph.generateId()
	} else {
		var vok, eok bool
		_, eok = graph.edges[string(id[:])]
		_, vok = graph.vertices[string(id[:])]
		if vok || eok {
			return nil, ErrAlreadyExists
		}
		idstr = id
	}
	graph.edgelock.Lock()
	defer graph.edgelock.Unlock()

	subject := outvertex.(*VertexMem)
	object := invertex.(*VertexMem)
	edge :=  NewEdgeMem(graph, idstr, subject, object, label)
	
	graph.edges[string(idstr[:])] = edge
	subject.addOutEdge(edge)
	object.addInEdge(edge)
	return edge, nil
}

func (graph *GraphMem) Edge(id []byte) (core.Edge, error) {
	if id == nil { return nil, ErrNilValue}
	graph.edgelock.RLock()
	defer graph.edgelock.RUnlock()
 	if val, ok := graph.edges[string(id[:])]; ok {
        return val, nil
    }
	return nil, nil
}

func (graph *GraphMem) DelEdge(edge core.Edge) error {
	if edge == nil { return ErrNilValue}
	//fmt.Printf("got to deledge\n")
	if _, ok := graph.edges[string(edge.Id()[:])]; !ok {
		return 	ErrDoesntExist
	}
	v, _ := edge.VertexOut()
	vertexout := v.(*VertexMem)
	v, _ = edge.VertexOut()
	vertexin := v.(*VertexMem)
	vertexout.delOutEdge(edge)
	vertexin.delOutEdge(edge)
	return nil
}

func (graph *GraphMem) Edges() ([]core.Edge, error) {
	graph.edgelock.RLock()
	defer graph.edgelock.RUnlock()
	edges := []core.Edge{}
	for _, e := range graph.edges {
		edges = append(edges, e)
	}
	return edges, nil
}

func (graph *GraphMem) EdgeCount() uint {
	return uint(len(graph.edges))
}

func (graph *GraphMem) VertexCount() uint {
	return uint(len(graph.vertices))
}

func (graph *GraphMem) Open(args ...interface{}) error {
	graph.vertices = make(map[string]*VertexMem)
	graph.edges = make(map[string]*EdgeMem)
	graph.vertexlock = &sync.RWMutex{}
	graph.edgelock = &sync.RWMutex{}
	graph.isopen = true
	return nil
}

func (graph *GraphMem) Close() error {
	graph.vertices = nil
	graph.edges = nil
	graph.vertexlock = nil
	graph.edgelock = nil
	graph.nextid = 1
	graph.isopen = false
	return nil
}

func (graph *GraphMem) IsOpen() bool {
	return graph.isopen
}

func (graph *GraphMem) generateId() []byte {
	id := ""
	//var vok bool
	for {
		id = strconv.FormatUint(graph.nextid, 16)
		graph.nextid++
		 _, vok := graph.vertices[id]
		 _, eok := graph.edges[id]
		if !vok  && !eok  {
			break
		}
	}
	return []byte(id)
}
