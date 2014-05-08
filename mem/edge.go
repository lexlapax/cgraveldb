package mem

import (
		"github.com/lexlapax/graveldb/core"
)

type EdgeMem struct {
	*AtomMem
	subject *VertexMem
	object *VertexMem
	label string
}

func NewEdgeMem(db *GraphMem, id []byte, subject *VertexMem, object *VertexMem, label string) *EdgeMem {
	edge := &EdgeMem{NewAtomMem(db, id, core.EdgeType), subject, object, label}
	return edge
}

func (edge *EdgeMem) Label() string {
	return edge.label
}

func (edge *EdgeMem) Vertex(direction core.Direction) (core.Vertex, error) {
	if direction == core.DirOut {
		return edge.subject, nil
	} else if direction == core.DirIn {
		return edge.object, nil
	} else {
		return nil, core.ErrDirAnyUnsupported
	}
}


func (edge *EdgeMem) VertexOut() (core.Vertex, error) {
	return edge.subject, nil
}

func (edge *EdgeMem) VertexIn() (core.Vertex, error) {
	return edge.object, nil
}
