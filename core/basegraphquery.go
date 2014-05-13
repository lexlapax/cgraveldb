package core

import (
		// "github.com/lexlapax/graveldb/util"
)

type BaseGraphQuery struct{
	*BaseQuery
	graph Graph
}

func NewBaseGraphQuery(graph Graph) *BaseGraphQuery {
	basequery := newBaseQuery()
	query := &BaseGraphQuery{basequery, graph}
	return query
}

func (query *BaseGraphQuery) Edges() []Edge {
	edges := []Edge{}
	if query == nil { return edges }
	return edges
}

func (query *BaseGraphQuery) IterEdges() <-chan Edge {
	ch := make(chan Edge)
	go func() {
		count := 0
		for edge := range query.graph.IterEdges() {
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


func (query *BaseGraphQuery) IterVertices()  <-chan Vertex {
	ch := make(chan Vertex)
	go func() {
		count := 0
		for vertex := range query.graph.IterVertices() {
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
