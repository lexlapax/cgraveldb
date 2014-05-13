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

func (query *BaseQuery) ConditionFilter(atom Atom) bool {
	meets = true
	

		for _, cond := range query.Conditions {
		switch cond.Has {
		case true:
			assert.True(t, cond.Value == nil)
			assert.True(t, cond.Has == true)
		case false:
			assert.True(t, cond.Value == nil)
			assert.True(t, cond.Has == false)
		}
	}
}

func (query *BaseVertexQuery) IterEdges() <-chan Edge {
	ch := make(chan Edge)
	go func() {
		count := 0
		for edge := range query.vertex.IterEdges(query.direction, query.labels...) {
			ch <- edge
			count++
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
			ch <- vertex
			count++
			if count == query.Max { break }
		}
		close(ch)
	}()
	return ch
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

