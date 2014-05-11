package core

import (
		// "github.com/lexlapax/graveldb/util"
)

type BaseGraphQuery struct{
	*BaseQuery
}

func NewBaseGraphQuery(vertex *Vertex) *BaseGraphQuery {
	basequery := newBaseQuery()
	query := &BaseGraphQuery{basequery}
	return query
}

func (query *BaseGraphQuery) Edges() []Edge {
	edges := []Edge{}
	if query == nil { return edges }
	return edges
}

func (query *BaseGraphQuery) Vertices() []Vertex {
	vertices := []Vertex{}
	if query == nil { return vertices }
	return vertices
}

