package core

import (
	"errors"
)

type Direction int
const (
	DirAny Direction = 0
	DirOut= 1
	DirIn = 2
)


type AtomType string
const (
	VertexType AtomType = "1"
	EdgeType ="2"
)

var (
	ErrDoesntExist = errors.New("the object with the id does not exist")
	ErrAlreadyExists = errors.New("the object with the id already exists")
	ErrNilValue = errors.New("value passed cannot be nil")
	ErrDirAnyUnsupported = errors.New("Any Direction not supported")
)

type GraphCaps interface {
	Persistent() bool
	SortedKeys() bool
	KeyIndex() bool
}

type Atom interface {
	Id() string
	Property(prop string) ([]byte, error)
	SetProperty(prop string, value []byte) error
	DelProperty(prop string) error
	PropertyKeys() ([]string, error)
}

type Edge interface {
	Atom
	Label() string
	Vertex(direction Direction) (Vertex, error)
	VertexOut() (Vertex, error)
	VertexIn() (Vertex, error)
	//String() string
}

// todo add channel interfaces for iteration
type Vertex interface {
	Atom
	IterEdges(direction Direction, labels ...string) <-chan Edge
	IterOutEdges(labels ...string) <-chan Edge
	IterInEdges(labels ...string) <-chan Edge
	Edges(direction Direction, labels ...string) ([]Edge, error)
	OutEdges(labels ...string) ([]Edge, error)
	InEdges(labels ...string) ([]Edge, error)
	IterVertices(direction Direction, labels ...string) <-chan Vertex
	Vertices(direction Direction, labels ...string) ([]Vertex, error)
	AddEdge(id string, invertex Vertex, label string) (Edge, error)
	//String() string
}

// todo add channel interfaces for iteration
type Graph interface {
	KeyIndexable
	Guid() string
	Capabilities() GraphCaps
	AddVertex(id string) (Vertex, error)
	Vertex(id string) (Vertex, error)
	DelVertex(vertex Vertex) error
	Vertices() ([]Vertex, error)
	AddEdge(id string, outvertex Vertex, invertex Vertex, label string) (Edge, error)
	Edge(id string) (Edge, error)
	DelEdge(edge Edge) error
	Edges() ([]Edge, error)
	IterVertices() <-chan Vertex
	IterEdges() <-chan Edge
	EdgeCount() uint
	VertexCount() uint
	IsOpen() bool
	Open(args ...interface{}) error
	Close() error
	Clear() error
}

type KeyIndexable interface {
	CreateKeyIndex(key string, atomType AtomType) error
	DropKeyIndex(key string, atomType AtomType) error
	IndexedKeys(atomType AtomType) []string
	//values are search for greedy match based on whitespace tokenization of values
	// ie.. allows for substring matches "abcd" matches "abcd efgh" , "abcd", "1234, abcd, efgh"
	VerticesWithProp(key string, value string) []Vertex
	EdgesWithProp(key string, value string) []Edge
}

type Query interface {
	HasKey (key string) Query
	NoKey (key string) Query
	HasKeyValue (key string, value []byte) Query
	NoKeyValue (key string, value []byte) Query
	Limit(limit int) Query
	IterEdges() <-chan Edge
	IterVertices() <-chan Vertex
}

type QueryGraph interface {
	Query	
}

type QueryVertex interface {
	Query
	HasLabels(labels ...string) QueryVertex
	Count() int
	Direction(dir Direction)
}

