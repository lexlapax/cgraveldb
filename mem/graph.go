package mem

import (
		"github.com/lexlapax/graveldb/core"
		"github.com/lexlapax/graveldb/util"
		"sync"
		"strconv"
		//"fmt"
)

const (
	GraphImpl                = "mem"
	)

var (
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
	caps *graphCaps
	keyindex *KeyIndex
	uuid string
}

func (graph *GraphMem) Capabilities() core.GraphCaps {
	return graph.caps
}

func (graph *GraphMem) AddVertex(id string) (core.Vertex, error) {
	graph.vertexlock.Lock()
	defer graph.vertexlock.Unlock()
	var idstr string
	if id == "" {
		idstr = graph.generateId()
	} else {
		var vok, eok bool
		_, eok = graph.edges[id]
		_, vok = graph.vertices[id]
		if vok || eok {
			return nil, core.ErrAlreadyExists
		}
		idstr = id
	}
	vertex := NewVertexMem(graph, idstr)
	graph.vertices[idstr] = vertex
	return vertex, nil
}

func (graph *GraphMem) Vertex(id string) (core.Vertex, error) {
	if id == "" { return nil, core.ErrNilValue}
	graph.vertexlock.RLock()
	defer graph.vertexlock.RUnlock()
 	if val, ok := graph.vertices[id]; ok {
        return val, nil
    }
	return nil, nil
}

func (graph *GraphMem) DelVertex(vertex core.Vertex) error {
	if vertex == nil { return core.ErrNilValue}
	graph.vertexlock.Lock()
	defer graph.vertexlock.Unlock()
	edges, _ := vertex.Edges(core.DirAny)
	//fmt.Printf("edges=%v\n",edges)
	for _, edge := range edges {
		graph.DelEdge(edge)
	}
	delete(graph.vertices, vertex.Id())
	return nil
}

func (graph *GraphMem) Vertices() ([]core.Vertex, error) {
	vertices := []core.Vertex{}
	for v := range graph.IterVertices() {
		vertices = append(vertices, v)
	}
	return vertices, nil
}

func (graph *GraphMem) IterVertices() <-chan core.Vertex {
	ch := make(chan core.Vertex)
	go func() {
		graph.vertexlock.RLock()
		defer graph.vertexlock.RUnlock()
		for _, v := range graph.vertices {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

func (graph *GraphMem) AddEdge(id string, outvertex core.Vertex, invertex core.Vertex, label string) (core.Edge, error) {
	if outvertex == nil || invertex == nil || outvertex.Id() == "" || invertex.Id() == "" {
		return nil, core.ErrNilValue
	}
	var idstr string
	if id == "" {
		idstr = graph.generateId()
	} else {
		var vok, eok bool
		_, eok = graph.edges[id]
		_, vok = graph.vertices[id]
		if vok || eok {
			return nil, core.ErrAlreadyExists
		}
		idstr = id
	}
	graph.edgelock.Lock()
	defer graph.edgelock.Unlock()

	subject := outvertex.(*VertexMem)
	object := invertex.(*VertexMem)
	edge :=  NewEdgeMem(graph, idstr, subject, object, label)
	
	graph.edges[idstr] = edge
	subject.addOutEdge(edge)
	object.addInEdge(edge)
	return edge, nil
}

func (graph *GraphMem) Edge(id string) (core.Edge, error) {
	if id == "" { return nil, core.ErrNilValue}
	graph.edgelock.RLock()
	defer graph.edgelock.RUnlock()
 	if val, ok := graph.edges[id]; ok {
        return val, nil
    }
	return nil, nil
}

func (graph *GraphMem) DelEdge(edge core.Edge) error {
	if edge == nil { return core.ErrNilValue}
	if _, ok := graph.edges[edge.Id()]; !ok {
		return 	core.ErrDoesntExist
	}
	graph.edgelock.Lock()
	defer graph.edgelock.Unlock()
	delete(graph.edges, edge.Id())
	v, _ := edge.VertexOut()
	vertexout := v.(*VertexMem)
	v, _ = edge.VertexIn()
	vertexin := v.(*VertexMem)
	vertexout.delOutEdge(edge)
	vertexin.delInEdge(edge)
	//fmt.Printf("got to deledge\n")
	return nil
}

func (graph *GraphMem) Edges() ([]core.Edge, error) {
	edges := []core.Edge{}
	for e := range graph.IterEdges() {
		edges = append(edges, e)
	}
	return edges, nil
}

func (graph *GraphMem) IterEdges() <-chan core.Edge {
	ch := make(chan core.Edge)
	go func() {
		graph.edgelock.RLock()
		defer graph.edgelock.RUnlock()
		for _, e := range graph.edges {
			ch <- e
		}
		close(ch)
	}()
	return ch
}


func (graph *GraphMem) Guid() string {
	return graph.uuid	
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
	graph.caps = new(graphCaps)
	graph.keyindex = NewKeyIndex()
	id, err := util.UUID()
	if err != nil { return err }
	graph.uuid = id
	return nil
}

func (graph *GraphMem) Close() error {
	graph.vertices = nil
	graph.edges = nil
	graph.vertexlock = nil
	graph.edgelock = nil
	graph.nextid = 1
	graph.isopen = false
	graph.keyindex.close()
	return nil
}

func (graph *GraphMem) Clear() error {
	graph.keyindex.clear()
	for _, v := range graph.vertices {
		graph.DelVertex(v)
	}
	return nil
}

func (graph *GraphMem) IsOpen() bool {
	return graph.isopen
}

func (graph *GraphMem) generateId() string {
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
	return id
}

func (graph *GraphMem) CreateKeyIndex(key string, atomType core.AtomType) error {
	return graph.keyindex.createKeyIndex(key, atomType)
}

func (graph *GraphMem) DropKeyIndex(key string, atomType core.AtomType) error {
	return graph.keyindex.dropKeyIndex(key, atomType)
}

func (graph *GraphMem) IndexedKeys(atomType core.AtomType) []string {
		return graph.keyindex.indexedKeys(atomType)
}

func (graph *GraphMem) VerticesWithProp(key string, value string) []core.Vertex {
	ids := graph.keyindex.searchIds(key, value, core.VertexType)
	vertices := []core.Vertex{}
	var vertex core.Vertex
	for _, idstring := range ids {
		vertex, _ = graph.Vertex(idstring)
		if vertex != nil {
			vertices = append(vertices, vertex)
		}
	}
	return vertices
}

func (graph *GraphMem) EdgesWithProp(key string, value string) []core.Edge {
	ids := graph.keyindex.searchIds(key, value, core.EdgeType)
	edges := []core.Edge{}
	var edge core.Edge
	for _, idstring := range ids {
		edge, _ = graph.Edge(idstring)
		if edge != nil {
			edges = append(edges, edge)
		}
	}
	return edges
}
