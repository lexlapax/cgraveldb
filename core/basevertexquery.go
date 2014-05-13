package core

import (
		// "github.com/lexlapax/graveldb/util"
)

type BaseVertexQuery struct{
	*BaseQuery
	vertex Vertex
	labels []string
	direction Direction
}

func NewBaseVertexQuery(vertex Vertex) *BaseVertexQuery {
	basequery := newBaseQuery()
	query := &BaseVertexQuery{basequery, vertex, []string{}, DirAny}
	return query
}


func (query *BaseVertexQuery) IterEdges() <-chan Edge {
	ch := make(chan Edge)
	go func() {
		count := 0
		for edge := range query.vertex.IterEdges(query.direction, query.labels...) {
			if query.ConditionFilter(edge) {
				ch <- edge
				count++
			}
			if count == query.Max { break }
		}
		close(ch)
	}()
	return ch
}


func (query *BaseVertexQuery) IterVertices()  <-chan Vertex {
	ch := make(chan Vertex)
	go func() {
		count := 0
		for vertex := range query.vertex.IterVertices(query.direction, query.labels...) {
			if query.ConditionFilter(vertex) {
				ch <- vertex
				count++
			}
			if count == query.Max { break }
		}
		close(ch)
	}()
	return ch
}

func (query *BaseVertexQuery) HasLabels(labels ...string) QueryVertex {
	for _, label := range labels {
		query.labels = append(query.labels, label)
	}
	return query
}

func (query *BaseVertexQuery) Direction(dir Direction) QueryVertex {
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

