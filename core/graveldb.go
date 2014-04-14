package core

import (
)

type Element interface {
	Property(prop string) []byte
	SetProperty(prop string, value []byte) (bool, error)
	DelProperty(prop string) []byte
	PropertyKeys() [][]byte
	Id() []byte
}

type Edge interface {
	Element
	Label() string
	VertexOut() Vertex
	VertexIn() Vertex
	String() string
}

type Vertex interface {
	Element
	Vertices() []Vertex
	OutEdges() []Edge
	InEdges() []Edge
	AddEdge(outvertex Vertex, label string) (Edge, error)
	String() string
}

type Graph interface {
	String() string
	AddVertex(id []byte) (Vertex, error)
//	Vertex(id []byte) Vertex
//	DelVertex(vertex Vertex) error
//	Vertices() []Vertex
//	AddEdge(id []byte, outvertex *Vertex, invertex *Vertex, label string) (*Edge, error)
//	Edge(id []byte) Edge
//	DelEdge(edge Edge) error
//	Edges() []Edge
//	EdgeCount() uint
//	VertexCount() uint
//	Close() (bool, error)
}