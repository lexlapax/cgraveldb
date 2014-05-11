package core

import (
		// "github.com/lexlapax/graveldb/util"
)

type BaseVertexQuery struct{
	*BaseQuery
	vertex *Vertex
	labels []string
	direction Direction
}

func NewBaseVertexQuery(vertex *Vertex) *BaseVertexQuery {
	basequery := newBaseQuery()
	query := &BaseVertexQuery{basequery, vertex, []string{}, DirAny}
	return query
}

// func (query *BaseVertexQuery) Edges() <-chan Edge {
// 	ch := make(chan Edge)
// 	go func() {
// 		for _, edge := range query.vertex.Edges(query.direction, query.labels) {

// 		}

// 	}()
// 	return ch
// }


func (query *BaseVertexQuery) Vertices() []Vertex {
	vertices := []Vertex{}
	if query == nil { return vertices }
	return vertices
}

func (query *BaseVertexQuery) HasLabels(labels ...string) *BaseVertexQuery {
	for _, label := range labels {
		query.labels = append(query.labels, label)
	}
	return query
}

func (query *BaseVertexQuery) Direction(dir Direction) *BaseVertexQuery {
	query.direction = dir
	return query
}


func (query *BaseVertexQuery) Count() int {
	return 0
}


// 	Edges() []Edge
// 	Vertices() []Vertex
// }

// type QueryGraph interface {
// 	Query	
// }

// type QueryVertex interface {
// 	Query
// 	HasLabels(labels ...string) BaseVertexQuery
// 	Count() uint
// }

