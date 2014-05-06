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
	Id() []byte
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
	Edges(direction Direction, labels ...string) ([]Edge, error)
	Vertices(direction Direction, labels ...string) ([]Vertex, error)
	OutEdges(labels ...string) ([]Edge, error)
	InEdges(labels ...string) ([]Edge, error)
	AddEdge(id []byte, invertex Vertex, label string) (Edge, error)
	//String() string
}

// todo add channel interfaces for iteration
type Graph interface {
	Capabilities() GraphCaps
	AddVertex(id []byte) (Vertex, error)
	Vertex(id []byte) (Vertex, error)
	DelVertex(vertex Vertex) error
	Vertices() ([]Vertex, error)
	AddEdge(id []byte, outvertex Vertex, invertex Vertex, label string) (Edge, error)
	Edge(id []byte) (Edge, error)
	DelEdge(edge Edge) error
	Edges() ([]Edge, error)
	EdgeCount() uint
	VertexCount() uint
	IsOpen() bool
	Open(args ...interface{}) error
	Close() error
	Clear() error
}

type KeyIndexableGraph interface {
	Graph
	CreateKeyIndex(key string, atomType AtomType) error
	DropKeyIndex(key string, atomType AtomType) error
	IndexedKeys(atomType AtomType) []string
}