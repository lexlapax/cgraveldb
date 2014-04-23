package mem

import (
		"github.com/lexlapax/graveldb/core"
		"errors"
)

var (
	ErrDirAnyUnsupported = errors.New("Any Direction not supported")
)

type EdgeMem struct {
	*NodeMem
	subject *VertexMem
	object *VertexMem
	label string
}

func NewEdgeMem(db *GraphMem, id []byte, subject *VertexMem, object *VertexMem, label string) *EdgeMem {
	edge := &EdgeMem{NewNodeMem(db, id, VertexType), subject, object, label}
	return edge
}

func (edge *EdgeMem) Label() string {
	return edge.label
}

func (edge *EdgeMem) Vertex(direction core.Direction) (core.Vertex, error) {
	if direction == core.DirForward {
		return edge.object, nil
	} else if direction == core.DirReverse {
		return edge.subject, nil
	} else {
		return nil, ErrDirAnyUnsupported
	}
}


func (edge *EdgeMem) VertexOut() (core.Vertex, error) {
	return edge.subject, nil
}

func (edge *EdgeMem) VertexIn() (core.Vertex, error) {
	return edge.object, nil
}
